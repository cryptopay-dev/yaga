package model

import (
	"fmt"

	"github.com/go-pg/pg/orm"
)

// Query for shadowing orm.Query
type Query = orm.Query

// Option closure
type Option func(*Query)

// Where is a shortcut for Where
func Where(where string, args ...interface{}) Option {
	return func(q *Query) {
		q.Where(where, args...)
	}
}

// WhereIn is a shortcut for Where and pg.In to work with IN operator:
func WhereIn(where string, args ...interface{}) Option {
	return func(q *Query) {
		q.WhereIn(where, args...)
	}
}

// WhereOr is a shortcut for Where and OR
func WhereOr(where string, args ...interface{}) Option {
	return func(q *Query) {
		q.WhereOr(where, args...)
	}
}

// Equal is shortcut for "field = ?", value
func Equal(field string, value interface{}) Option {
	return func(q *Query) {
		q.Where(fmt.Sprintf("%s = ?", field), value)
	}
}
