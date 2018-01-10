package orm

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/go-pg/pg/internal"
	"github.com/go-pg/pg/types"
)

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()
var ipType = reflect.TypeOf((*net.IP)(nil)).Elem()
var ipNetType = reflect.TypeOf((*net.IPNet)(nil)).Elem()
var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
var nullBoolType = reflect.TypeOf((*sql.NullBool)(nil)).Elem()
var nullFloatType = reflect.TypeOf((*sql.NullFloat64)(nil)).Elem()
var nullIntType = reflect.TypeOf((*sql.NullInt64)(nil)).Elem()
var nullStringType = reflect.TypeOf((*sql.NullString)(nil)).Elem()

type Table struct {
	Type       reflect.Type
	zeroStruct reflect.Value

	TypeName  string
	Name      types.Q
	Alias     types.Q
	ModelName string

	Fields     []*Field // PKs + DataFields
	PKs        []*Field
	DataFields []*Field
	FieldsMap  map[string]*Field

	Methods   map[string]*Method
	Relations map[string]*Relation

	flags uint8
}

func (t *Table) SetFlag(flag uint8) {
	t.flags |= flag
}

func (t *Table) HasFlag(flag uint8) bool {
	if t == nil {
		return false
	}
	return t.flags&flag != 0
}

func (t *Table) HasField(field string) bool {
	_, err := t.GetField(field)
	return err == nil
}

func (t *Table) checkPKs() error {
	if len(t.PKs) == 0 {
		return fmt.Errorf("model=%s does not have primary keys", t.Type.Name())
	}
	return nil
}

func (t *Table) AddField(field *Field) {
	t.Fields = append(t.Fields, field)
	if field.HasFlag(PrimaryKeyFlag) {
		t.PKs = append(t.PKs, field)
	} else {
		t.DataFields = append(t.DataFields, field)
	}
	t.FieldsMap[field.SQLName] = field
}

func (t *Table) RemoveField(field *Field) {
	t.Fields = removeField(t.Fields, field)
	if field.HasFlag(PrimaryKeyFlag) {
		t.PKs = removeField(t.PKs, field)
	} else {
		t.DataFields = removeField(t.DataFields, field)
	}
	delete(t.FieldsMap, field.SQLName)
}

func removeField(fields []*Field, field *Field) []*Field {
	for i, f := range fields {
		if f == field {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}
	return fields
}

func (t *Table) GetField(fieldName string) (*Field, error) {
	field, ok := t.FieldsMap[fieldName]
	if !ok {
		return nil, fmt.Errorf("can't find column=%s in table=%s", fieldName, t.Name)
	}
	return field, nil
}

func (t *Table) AppendParam(b []byte, strct reflect.Value, name string) ([]byte, bool) {
	if field, ok := t.FieldsMap[name]; ok {
		b = field.AppendValue(b, strct, 1)
		return b, true
	}

	if method, ok := t.Methods[name]; ok {
		b = method.AppendValue(b, strct.Addr(), 1)
		return b, true
	}

	return b, false
}

func (t *Table) addRelation(rel *Relation) {
	if t.Relations == nil {
		t.Relations = make(map[string]*Relation)
	}
	t.Relations[rel.Field.GoName] = rel
}

func newTable(typ reflect.Type) *Table {
	table, ok := Tables.tables[typ]
	if ok {
		return table
	}

	modelName := internal.Underscore(typ.Name())
	table = &Table{
		Type:       typ,
		zeroStruct: reflect.Zero(typ),

		TypeName:  internal.ToExported(typ.Name()),
		Name:      types.Q(types.AppendField(nil, tableNameInflector(modelName), 1)),
		Alias:     types.Q(types.AppendField(nil, modelName, 1)),
		ModelName: modelName,

		Fields:    make([]*Field, 0, typ.NumField()),
		FieldsMap: make(map[string]*Field, typ.NumField()),
	}
	Tables.tables[typ] = table

	table.addFields(typ, nil)
	typ = reflect.PtrTo(typ)

	if typ.Implements(afterQueryHookType) {
		table.SetFlag(AfterQueryHookFlag)
	}
	if typ.Implements(afterSelectHookType) {
		table.SetFlag(AfterSelectHookFlag)
	}
	if typ.Implements(beforeInsertHookType) {
		table.SetFlag(BeforeInsertHookFlag)
	}
	if typ.Implements(afterInsertHookType) {
		table.SetFlag(AfterInsertHookFlag)
	}
	if typ.Implements(beforeUpdateHookType) {
		table.SetFlag(BeforeUpdateHookFlag)
	}
	if typ.Implements(afterUpdateHookType) {
		table.SetFlag(AfterUpdateHookFlag)
	}
	if typ.Implements(beforeDeleteHookType) {
		table.SetFlag(BeforeDeleteHookFlag)
	}
	if typ.Implements(afterDeleteHookType) {
		table.SetFlag(AfterDeleteHookFlag)
	}

	if table.Methods == nil {
		table.Methods = make(map[string]*Method)
	}
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.PkgPath != "" {
			continue
		}
		if m.Type.NumIn() > 1 {
			continue
		}
		if m.Type.NumOut() != 1 {
			continue
		}

		retType := m.Type.Out(0)
		method := Method{
			Index: m.Index,

			appender: types.Appender(retType),
		}

		table.Methods[m.Name] = &method
	}

	return table
}

func (t *Table) addFields(typ reflect.Type, baseIndex []int) {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		// Make a copy so slice is not shared between fields.
		var index []int
		index = append(index, baseIndex...)

		if f.Anonymous {
			sqlTag := f.Tag.Get("sql")
			if sqlTag == "-" {
				continue
			}

			embeddedTable := newTable(indirectType(f.Type))

			pgTag := parseTag(f.Tag.Get("pg"))
			if _, ok := pgTag.Options["override"]; ok {
				t.TypeName = embeddedTable.TypeName
				t.Name = embeddedTable.Name
				t.Alias = embeddedTable.Alias
				t.ModelName = embeddedTable.ModelName
			}

			t.addFields(embeddedTable.Type, append(index, f.Index...))
			continue
		}

		field := t.newField(f, index)
		if field != nil {
			t.AddField(field)
		}
	}
}

func (t *Table) getField(name string) *Field {
	for _, f := range t.Fields {
		if f.GoName == name {
			return f
		}
	}

	f, ok := t.Type.FieldByName(name)
	if !ok {
		return nil
	}
	return t.newField(f, nil)
}

func (t *Table) newField(f reflect.StructField, index []int) *Field {
	sqlTag := parseTag(f.Tag.Get("sql"))

	switch f.Name {
	case "tableName", "TableName":
		if index != nil {
			return nil
		}
		if sqlTag.Name != "" {
			if isPostgresKeyword(sqlTag.Name) {
				sqlTag.Name = `"` + sqlTag.Name + `"`
			}
			t.Name = types.Q(sqlTag.Name)
		}
		if alias, ok := sqlTag.Options["alias"]; ok {
			t.Alias = types.Q(alias)
		}
		return nil
	}

	if f.PkgPath != "" {
		return nil
	}

	skip := sqlTag.Name == "-"
	if skip || sqlTag.Name == "" {
		sqlTag.Name = internal.Underscore(f.Name)
	}

	index = append(index, f.Index...)
	if field, ok := t.FieldsMap[sqlTag.Name]; ok {
		if indexEqual(field.Index, index) {
			return field
		}
		t.RemoveField(field)
	}

	field := Field{
		Type: indirectType(f.Type),

		GoName:  f.Name,
		SQLName: sqlTag.Name,
		Column:  types.Q(types.AppendField(nil, sqlTag.Name, 1)),

		Index: index,
	}

	if _, ok := sqlTag.Options["notnull"]; ok {
		field.SetFlag(NotNullFlag)
	}
	if _, ok := sqlTag.Options["unique"]; ok {
		field.SetFlag(UniqueFlag)
	}
	if v, ok := sqlTag.Options["default"]; ok {
		v, ok = unquote(v)
		if ok {
			field.Default = types.Q(types.AppendString(nil, v, 1))
		} else {
			field.Default = types.Q(v)
		}
	}

	if len(t.PKs) == 0 && (field.SQLName == "id" || field.SQLName == "uuid") {
		field.SetFlag(PrimaryKeyFlag)
	} else if _, ok := sqlTag.Options["pk"]; ok {
		field.SetFlag(PrimaryKeyFlag)
	} else if strings.HasSuffix(field.SQLName, "_id") ||
		strings.HasSuffix(field.SQLName, "_uuid") {
		field.SetFlag(ForeignKeyFlag)
	}

	pgTag := parseTag(f.Tag.Get("pg"))
	if _, ok := pgTag.Options["array"]; ok {
		field.SetFlag(ArrayFlag)
	}

	field.SQLType = fieldSQLType(&field, sqlTag)
	if strings.HasSuffix(field.SQLType, "[]") {
		field.SetFlag(ArrayFlag)
	}

	if _, ok := pgTag.Options["json_use_number"]; ok {
		field.append = types.Appender(f.Type)
		field.scan = scanJSONValue
	} else if field.HasFlag(ArrayFlag) {
		field.append = types.ArrayAppender(f.Type)
		field.scan = types.ArrayScanner(f.Type)
	} else if _, ok := pgTag.Options["hstore"]; ok {
		field.append = types.HstoreAppender(f.Type)
		field.scan = types.HstoreScanner(f.Type)
	} else {
		field.append = types.Appender(f.Type)
		field.scan = types.Scanner(f.Type)
	}
	field.isZero = isZeroFunc(f.Type)

	if !skip && isColumn(f.Type) {
		return &field
	}

	switch field.Type.Kind() {
	case reflect.Slice:
		elemType := indirectType(field.Type.Elem())
		if elemType.Kind() != reflect.Struct {
			break
		}

		joinTable := newTable(elemType)

		fk, ok := pgTag.Options["fk"]
		if !ok {
			fk = t.TypeName
		}

		if m2mTable, _ := pgTag.Options["many2many"]; m2mTable != "" {
			m2mTableAlias := m2mTable
			if ind := strings.IndexByte(m2mTable, '.'); ind >= 0 {
				m2mTableAlias = m2mTable[ind+1:]
			}

			joinFK, ok := pgTag.Options["joinFK"]
			if !ok {
				joinFK = joinTable.TypeName
			}

			t.addRelation(&Relation{
				Type:          Many2ManyRelation,
				Field:         &field,
				JoinTable:     joinTable,
				M2MTableName:  types.Q(m2mTable),
				M2MTableAlias: types.Q(m2mTableAlias),
				BasePrefix:    internal.Underscore(fk + "_"),
				JoinPrefix:    internal.Underscore(joinFK + "_"),
			})
			return nil
		}

		s, polymorphic := pgTag.Options["polymorphic"]
		if polymorphic {
			fk = s
		}

		fks := foreignKeys(t, joinTable, fk, t.TypeName)
		if len(fks) > 0 {
			t.addRelation(&Relation{
				Type:        HasManyRelation,
				Polymorphic: polymorphic,
				Field:       &field,
				FKs:         fks,
				JoinTable:   joinTable,
				BasePrefix:  internal.Underscore(fk + "_"),
			})
			return nil
		}
	case reflect.Struct:
		joinTable := newTable(field.Type)
		if len(joinTable.Fields) == 0 {
			break
		}

		for _, ff := range joinTable.FieldsMap {
			ff = ff.Copy()
			ff.SQLName = field.SQLName + "__" + ff.SQLName
			ff.Column = types.Q(types.AppendField(nil, ff.SQLName, 1))
			ff.Index = append(field.Index, ff.Index...)
			if _, ok := t.FieldsMap[ff.SQLName]; !ok {
				t.FieldsMap[ff.SQLName] = ff
			}
		}

		if t.tryHasOne(joinTable, &field, pgTag) ||
			t.tryBelongsToOne(joinTable, &field, pgTag) {
			t.FieldsMap[field.SQLName] = &field
			return nil
		}
	}

	if skip {
		t.FieldsMap[field.SQLName] = &field
		return nil
	}

	return &field
}

func isPostgresKeyword(s string) bool {
	switch s {
	case "user":
		return true
	}
	return false
}

func isColumn(typ reflect.Type) bool {
	return typ.Implements(scannerType) || reflect.PtrTo(typ).Implements(scannerType)
}

func fieldSQLType(field *Field, sqlTag *tag) string {
	if typ, ok := sqlTag.Options["type"]; ok {
		field.SetFlag(customTypeFlag)
		typ, _ := unquote(typ)
		return typ
	}

	if field.HasFlag(ArrayFlag) {
		sqlType := sqlType(field.Type.Elem())
		return sqlType + "[]"
	}

	sqlType := sqlType(field.Type)
	if field.HasFlag(PrimaryKeyFlag) {
		switch sqlType {
		case "smallint":
			return "smallserial"
		case "integer":
			return "serial"
		case "bigint":
			return "bigserial"
		}
	}

	switch sqlType {
	case "timestamptz":
		field.SetFlag(customTypeFlag)
	}

	return sqlType
}

func sqlType(typ reflect.Type) string {
	switch typ {
	case timeType:
		return "timestamptz"
	case ipType:
		return "inet"
	case ipNetType:
		return "cidr"
	case nullBoolType:
		return "boolean"
	case nullFloatType:
		return "double precision"
	case nullIntType:
		return "bigint"
	case nullStringType:
		return "text"
	}

	switch typ.Kind() {
	case reflect.Int8, reflect.Uint8, reflect.Int16:
		return "smallint"
	case reflect.Uint16, reflect.Int32:
		return "integer"
	case reflect.Uint32, reflect.Int64, reflect.Int:
		return "bigint"
	case reflect.Uint, reflect.Uint64:
		return "decimal"
	case reflect.Float32:
		return "real"
	case reflect.Float64:
		return "double precision"
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "text"
	case reflect.Map, reflect.Struct:
		return "jsonb"
	case reflect.Array, reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return "bytea"
		}
		return "jsonb"
	default:
		return typ.Kind().String()
	}
}

func (t *Table) tryHasOne(joinTable *Table, field *Field, tag *tag) bool {
	fk, ok := tag.Options["fk"]
	if !ok {
		fk = field.GoName
	}

	fks := foreignKeys(joinTable, t, fk, field.GoName)
	if len(fks) > 0 {
		t.addRelation(&Relation{
			Type:      HasOneRelation,
			Field:     field,
			FKs:       fks,
			JoinTable: joinTable,
		})
		return true
	}
	return false
}

func (t *Table) tryBelongsToOne(joinTable *Table, field *Field, tag *tag) bool {
	fk, ok := tag.Options["fk"]
	if !ok {
		fk = t.TypeName
	}

	fks := foreignKeys(t, joinTable, fk, t.TypeName)
	if len(fks) > 0 {
		t.addRelation(&Relation{
			Type:      BelongsToRelation,
			Field:     field,
			FKs:       fks,
			JoinTable: joinTable,
		})
		return true
	}
	return false
}

func foreignKeys(base, join *Table, fk, fieldName string) []*Field {
	var fks []*Field

	for _, pk := range base.PKs {
		fkName := fk + pk.GoName
		if f := join.getField(fkName); f != nil {
			fks = append(fks, f)
		}
	}

	if len(fks) > 0 {
		return fks
	}

	if fk == "" {
		return nil
	}

	if fk != fieldName {
		f := join.getField(fk)
		if f != nil {
			fks = append(fks, f)
			return fks
		}
	}

	for _, suffix := range []string{"Id", "ID", "UUID"} {
		f := join.getField(fk + suffix)
		if f != nil {
			fks = append(fks, f)
			return fks
		}
	}

	return nil
}

func scanJSONValue(v reflect.Value, b []byte) error {
	if !v.CanSet() {
		return fmt.Errorf("pg: Scan(non-pointer %s)", v.Type())
	}
	if b == nil {
		v.Set(reflect.New(v.Type()).Elem())
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(v.Addr().Interface())
}
