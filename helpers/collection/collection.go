package collection

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/go-pg/pg/orm"
)

const defaultOrdersLimit = 25

// Filter interface
type Filter interface {
	Pager() Paginate
	Apply(*orm.Query) error
}

// Paginate is part of filter structure
type Paginate struct {
	Offset int `query:"offset" form:"offset" json:"offset"`
	Limit  int `query:"limit" form:"limit" json:"limit"`
}

// Pager returns Paginate-struct
func (p Paginate) Pager() Paginate {
	var limit = defaultOrdersLimit

	if p.Limit > 0 {
		limit = p.Limit
	}

	return Paginate{
		Offset: p.Offset,
		Limit:  limit,
	}
}

// Collection response answer
type Collection struct {
	Total  int   `json:"total"`
	Offset int   `json:"offset"`
	Items  Items `json:"items"`
}

// ModelsFetcher closure to fetch and format Items
type ModelsFetcher = func(*orm.Query) (Items, error)

// Items is formatted response of models
type Items = []interface{}

// Options to Format list-response
type Options struct {
	Query   *orm.Query
	Fetcher ModelsFetcher
	Filter  Filter
}

// Format collection as response-answer (Collection)
func Format(ctx web.Context, opts Options) error {
	if opts.Query == nil || opts.Filter == nil || opts.Fetcher == nil {
		return nil
	}

	var (
		err      error
		pager    = opts.Filter.Pager()
		response Collection
	)

	if err = opts.Filter.Apply(opts.Query); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	if response.Total, err = opts.Query.Count(); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusInternalServerError, err.Error())
	}

	// Pagination:
	opts.Query.Limit(pager.Limit)
	opts.Query.Offset(pager.Offset)

	response.Offset = pager.Offset

	if response.Items, err = opts.Fetcher(opts.Query); err != nil {
		ctx.Logger().Error(err.Error())

		if logicError, ok := err.(*errors.LogicError); ok {
			return logicError
		}

		return errors.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}
