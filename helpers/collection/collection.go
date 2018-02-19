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
	// Former interface
	Former interface {
		ApplyFilter(*Options) error
		ApplyPager(*Options) error
		ApplySorter(*Options) error
	}

	// Fetcher closure to fetch and format Items
	Fetcher = func(*Options) (Items, error)

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
		Fetcher Fetcher
		Former  Former
	}
)

// Response collection as response-answer (Collection)
//
// if you do not want to use opts.Query
// then see https://gist.github.com/jenchik/05a6a1bc80e7199203b8fa122cdaf922
func Response(ctx web.Context, opts Options) error {
	if opts.Query == nil || opts.Former == nil || opts.Fetcher == nil {
		return nil
	}

	var (
		err      error
		response Collection
	)

	if err = opts.Former.ApplyFilter(&opts); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	if response.Total, err = opts.Query.Count(); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusInternalServerError, err.Error())
	}

	if err = opts.Former.ApplyPager(&opts); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	if err = opts.Former.ApplySorter(&opts); err != nil {
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
