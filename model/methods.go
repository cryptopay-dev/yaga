package model

import (
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
)

func Create(db orm.DB, v interface{}) (int, error) {
	res, err := db.Model(v).Insert()
	if err != nil {
		return 0, errors.Wrap(err, "Error inserting record")
	}

	return res.RowsAffected(), nil
}

func Delete(db orm.DB, v interface{}) (int, error) {
	res, err := db.Model(v).Delete()
	if err != nil {
		return 0, errors.Wrap(err, "Error deleting record")
	}

	return res.RowsAffected(), nil
}

func Update(db orm.DB, v interface{}, column ...string) (int, error) {
	res, err := db.Model(v).Column(column...).Update()
	if err != nil {
		return 0, errors.Wrap(err, "Error updating record")
	}

	return res.RowsAffected(), nil
}

type Conditions map[string]interface{}

func queryFilter(db orm.DB, filter Conditions, v interface{}) *orm.Query {
	q := db.Model(v)
	for name, value := range filter {
		q.Where(name+"=?", value)
	}

	return q
}

func Find(db orm.DB, filter Conditions, v interface{}) error {
	f := queryFilter(db, filter, v)

	if err := f.Select(); err != nil {
		return errors.Wrap(err, "Error finding records")
	}

	return nil
}

func FindOneByID(db orm.DB, id int64, v interface{}) error {
	return FindOne(db, Conditions{"id": id}, v)
}

func FindOne(db orm.DB, filter Conditions, v interface{}) error {
	f := queryFilter(db, filter, v)

	if err := f.First(); err != nil {
		return errors.Wrap(err, "Error finding record")
	}

	return nil
}

func FindOneForUpdate(db orm.DB, filter Conditions, v interface{}) error {
	f := queryFilter(db, filter, v)

	if err := f.For("UPDATE").First(); err != nil {
		return errors.Wrap(err, "Error finding record")
	}

	return nil
}
