package model

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// ErrNoRows in database
var ErrNoRows = pg.ErrNoRows

// Create row in database
func Create(db orm.DB, model interface{}) (int, error) {
	res, err := db.Model(model).Insert()
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}

// Delete row from database
func Delete(db orm.DB, model interface{}, opts ...Option) (int, error) {
	res, err := queryFilter(db, model, opts...).Delete()
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}

// Update row in database
func Update(db orm.DB, v interface{}, column ...string) (int, error) {
	res, err := db.Model(v).Column(column...).Update()
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}

func queryFilter(db orm.DB, model interface{}, opts ...Option) *orm.Query {
	q := db.Model(model)
	for _, o := range opts {
		o(q)
	}

	return q
}

// Find row in database
func Find(db orm.DB, model interface{}, opts ...Option) error {
	return queryFilter(db, model, opts...).Select()
}

// FindByID row in database
func FindByID(db orm.DB, model interface{}, id interface{}) error {
	return FindOne(db, model, Equal("id", id))
}

// FindOne row in database
func FindOne(db orm.DB, model interface{}, opts ...Option) error {
	return queryFilter(db, model, opts...).First()
}

// Exists for check for row in database
func Exists(db orm.DB, model interface{}, opts ...Option) (bool, error) {
	n, err := queryFilter(db, model, opts...).Count()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

// FindOneForUpdate row in database and row-level locking
func FindOneForUpdate(db orm.DB, model interface{}, opts ...Option) error {
	return queryFilter(db, model, opts...).For("UPDATE").First()
}
