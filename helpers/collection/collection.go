package collection

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/go-pg/pg/orm"
)

// DefaultItemsLimit per page limit
var DefaultItemsLimit = 50

// SeparatorOrder
var SeparatorOrder = ","

type (
	// Filter interface
	Filter interface {
		ApplyFilter(*Options, *Collection) error
		ApplyPager(*Options, *Collection) error
		ApplySorter(*Options, *Collection) error
	}

	// ModelsFetcher closure to fetch and format Items
	ModelsFetcher = func(*Options) (Items, error)

	// Items is formatted response of models
	Items = []interface{}

	// Collection response answer
	Collection struct {
		Total  int   `json:"total"`
		Offset int   `json:"offset"`
		Items  Items `json:"items"`
	}

	// Options to Format list-response
	Options struct {
		Query   *orm.Query
		Fetcher ModelsFetcher
		Filter  Filter
	}
)

// FormatQuery collection as response-answer (Collection)
func FormatQuery(ctx web.Context, opts Options) error {
	if opts.Query == nil || opts.Filter == nil || opts.Fetcher == nil {
		return nil
	}

	var (
		err      error
		response Collection
	)

	if err = opts.Filter.ApplyFilter(&opts, &response); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	if response.Total, err = opts.Query.Count(); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusInternalServerError, err.Error())
	}

	if err = opts.Filter.ApplyPager(&opts, &response); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	if err = opts.Filter.ApplySorter(&opts, &response); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	return format(ctx, &opts, response)
}

// FormatSimple collection as response-answer (Collection)
func FormatSimple(ctx web.Context, opts Options) error {
	if opts.Filter == nil || opts.Fetcher == nil {
		return nil
	}

	var (
		err      error
		response Collection
	)

	if err = opts.Filter.ApplyFilter(&opts, &response); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	if err = opts.Filter.ApplyPager(&opts, &response); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	return format(ctx, &opts, response)
}

func format(ctx web.Context, opts *Options, response Collection) (err error) {
	if response.Items, err = opts.Fetcher(opts); err != nil {
		ctx.Logger().Error(err.Error())

		if logicError, ok := err.(*errors.LogicError); ok {
			return logicError
		}

		return errors.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}
