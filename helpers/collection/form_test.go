package collection

import (
	"strings"
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/stretchr/testify/assert"
)

var _ sorter = &mockQuery{}

type mockQuery struct {
	order string
}

func (q *mockQuery) Order(orders ...string) *orm.Query {
	q.order = strings.Join(orders, ",")

	return nil
}

var fieldsForTestSorter = []struct {
	in     string
	result string
}{
	// invalid
	{"", ""},
	{"a", ""},
	{"abc", ""},
	{"a,bc", ""},
	{"a,bc,def", ""},
	{"a,bc,de,+f", ""},
	{"a,bc,de,-f", ""},
	{"+,bc,de,-f", ""},
	{"a,+,de,-f", ""},
	{"a,bc,-,+f", ""},
	{"a,+,-,+f", ""},
	{"a,+-,+", "- ASC"},

	// ASC
	{"+a", "a ASC"},
	{"+abc", "abc ASC"},
	{"+a,bc", "a ASC"},
	{"a,+bc", "bc ASC"},
	{"ab,+c", "c ASC"},

	// DESC
	{"-a", "a DESC"},
	{"-abc", "abc DESC"},
	{"-a,bc", "a DESC"},
	{"a,-bc", "bc DESC"},
	{"ab,-c", "c DESC"},

	// both
	{"+a,-bc", "a ASC,bc DESC"},
	{"-a,+bc", "a DESC,bc ASC"},
	{"+a,+bc,-def", "a ASC,bc ASC,def DESC"},
	{"-a,+bc,-def", "a DESC,bc ASC,def DESC"},
	{"+a,-bc,-def", "a ASC,bc DESC,def DESC"},
	{"+a,-bc,+def", "a ASC,bc DESC,def ASC"},

	// both and invalid
	{"+a,+bc,def", "a ASC,bc ASC"},
	{"+a,-bc,def", "a ASC,bc DESC"},
	{"-a,+bc,def", "a DESC,bc ASC"},
	{"-a,-bc,def", "a DESC,bc DESC"},
	{"+a,bc,+def", "a ASC,def ASC"},
	{"+a,bc,-def", "a ASC,def DESC"},
	{"-a,bc,+def", "a DESC,def ASC"},
	{"-a,bc,-def", "a DESC,def DESC"},

	// > 3
	{"+a,+bc,de,+f", "a ASC,bc ASC"},
	{"+a,-bc,de,-f", "a ASC,bc DESC"},
	{"a,+bc,-de,+f", "bc ASC,de DESC"},
	{"a,-bc,-de,-f", "bc DESC,de DESC"},

	// case insensetive
	{"+A,+bc,-def", "a ASC,bc ASC,def DESC"},
	{"-a,+BC,-def", "a DESC,bc ASC,def DESC"},
	{"+a,-bc,-DEF", "a ASC,bc DESC,def DESC"},
	{"+A,-bc,+DEF", "a ASC,bc DESC,def ASC"},

	// repeats
	{"+bc,-bc,+bc", "bc ASC"},
	{"-bc,+bc,+bc", "bc ASC"},
	{"+bc,+bc,-bc", "bc DESC"},
	{"+a,+bc,+BC", "a ASC,bc ASC"},
	{"+a,-bc,+BC", "a ASC,bc ASC"},
	{"-A,+bc,-a", "a DESC,bc ASC"},
	{"-A,-bc,+a", "a ASC,bc DESC"},
}

func TestSorter(t *testing.T) {
	query := new(mockQuery)

	for _, item := range fieldsForTestSorter {
		applySorter(item.in, query)
		assert.Equal(t, item.result, query.order, "in: "+item.in)
		query.order = ""
	}
}
