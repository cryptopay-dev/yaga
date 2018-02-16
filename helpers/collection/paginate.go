package collection

import (
	"encoding/json"
	"strings"

	"github.com/go-pg/pg/orm"
)

// Paginate is part of filter structure
type Paginate struct {
	Order  string      `query:"order" form:"order" json:"order"`
	Offset json.Number `query:"offset" form:"offset" json:"offset" validate:"omitempty,gte=0"`
	Limit  json.Number `query:"limit" form:"limit" json:"limit" validate:"omitempty,gte=0"`
}

type sorter interface {
	Order(...string) *orm.Query
}

type order struct {
	d byte
	i int
}

func applySorter(inOrder string, query sorter) {
	var (
		found bool
		idx   int
		o     order
	)

	slice := strings.SplitN(strings.ToLower(inOrder), SeparatorOrder, 4)
	m := make(map[string]order, len(slice))
	for i, field := range slice {
		if i >= 3 {
			// sorting a maximum of 3 fields
			break
		}

		if len(field) < 2 || (field[0] != '-' && field[0] != '+') {
			continue
		}

		if o, found = m[field[1:]]; found {
			m[field[1:]] = order{field[0], o.i}
		} else {
			m[field[1:]] = order{field[0], idx}
			idx++
		}
	}

	for field, o := range m {
		if o.d == '-' {
			slice[o.i] = field + " DESC"
		} else {
			slice[o.i] = field + " ASC"
		}
	}

	if len(m) > 0 {
		query.Order(slice[:len(m)]...)
	}
}

// ApplySorter your filters to a Query
func (p *Paginate) ApplySorter(opts *Options, res *Collection) error {
	if len(p.Order) == 0 || opts.Query == nil {
		return nil
	}

	applySorter(p.Order, opts.Query)

	return nil
}

// ApplyPager
func (p *Paginate) ApplyPager(opts *Options, res *Collection) (err error) {
	var (
		limit = DefaultItemsLimit
		val   int64
	)

	if val, err = p.Offset.Int64(); err != nil {
		return err
	}
	res.Offset = int(val)

	if val, err = p.Limit.Int64(); err != nil {
		return err
	} else if val > 0 {
		limit = int(val)
	}

	if opts.Query != nil {
		opts.Query.Limit(limit)
		opts.Query.Offset(res.Offset)
	}

	return nil
}

// ApplyFilter your filters to a Query
func (*Paginate) ApplyFilter(*Options, *Collection) (err error) {
	// ...

	return nil
}
