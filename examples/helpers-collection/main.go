package main

import (
	"net/http"
	"os"
	"time"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/helpers/collection"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func main() {
	log := nop.New()

	e := web.New(web.Options{
		Logger: log,
	})

	ctrl := Controller{
		DB: pg.Connect(&pg.Options{
			Addr:     os.Getenv("DATABASE_ADDR"),
			User:     os.Getenv("DATABASE_USER"),
			Database: os.Getenv("DATABASE_DATABASE"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			PoolSize: 2,
		}),
	}

	e.GET("/", ctrl.ListCollections)

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}

}

// Controller struct
type Controller struct {
	DB *pg.DB
}

// MyModel struct of PG-table
type MyModel struct {
	ID        int64 `sql:"-,pk"`
	SomeField string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ModelResponse struct for response
type ModelResponse struct {
	ID        int64     `json:"id"`
	SomeField string    `json:"some_field"`
	CreatedAt time.Time `json:"created_at"`
}

// FormFilter with embedded Pager
type FormFilter struct {
	collection.Pager
}

// Apply your filters to a Query
func (f FormFilter) Apply(query *orm.Query) error {
	// ...

	return nil
}

func formatModelResponse(model *MyModel) *ModelResponse {
	if model == nil {
		return nil
	}

	return &ModelResponse{
		ID:        model.ID,
		SomeField: model.SomeField,
		CreatedAt: model.CreatedAt,
	}
}

// ListCollections handler
func (c *Controller) ListCollections(ctx web.Context) error {
	var req FormFilter

	if err := ctx.Bind(&req); err != nil {
		ctx.Logger().Error(err.Error())
		return errors.NewError(http.StatusBadRequest, err.Error())
	}

	return collection.Format(ctx, collection.Options{
		Query:  c.DB.Model(&MyModel{}),
		Filter: req,
		Fetcher: func(query *orm.Query) (collection.Items, error) {
			models := make([]MyModel, 0)

			if err := query.Select(&models); err != nil {
				return nil, errors.NewError(http.StatusInternalServerError, err.Error())
			}

			items := make(collection.Items, 0, len(models))

			for _, model := range models {
				items = append(items, formatModelResponse(&model))
			}

			return items, nil
		},
	})
}
