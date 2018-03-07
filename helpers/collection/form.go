package collection

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-pg/pg/orm"
)

const errNumericField = "'%s' must be a number"

// Form is part of filter structure
type Form struct {
	Order  string      `query:"order" form:"order" json:"order"`
	Offset json.Number `query:"offset" form:"offset" json:"offset" validate:"omitempty,numeric"`
	Limit  json.Number `query:"limit" form:"limit" json:"limit" validate:"omitempty,numeric"`
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

// ApplySorter your sorter to a Query
func (f *Form) ApplySorter(opts *Options) error {
	if len(f.Order) == 0 || opts.Query == nil {
		return nil
	}

	applySorter(f.Order, opts.Query)

	return nil
}

// ApplyPager to the opts.Query
func (f *Form) ApplyPager(opts *Options) (err error) {
	if opts.Query == nil {
		return nil
	}

	var (
		limit = DefaultItemsLimit
		val   int64
	)

	if len(f.Limit) > 0 {
		if val, err = f.Limit.Int64(); err != nil {
			return fmt.Errorf(errNumericField, "limit")
		} else if val <= 0 {
			return fmt.Errorf("'limit' must be greater than '0'")
		}
		if val > 0 {
			limit = int(val)
		}
	}
	opts.Query.Limit(limit)

	if len(f.Offset) > 0 {
		if val, err = f.Offset.Int64(); err != nil {
			return fmt.Errorf(errNumericField, "offset")
		} else if val < 0 {
			return fmt.Errorf("'offset' must be greater or equal than '0'")
		}
		opts.Query.Offset(int(val))
	}

	return nil
}

// ApplyFilter your filters to a Query
func (*Form) ApplyFilter(*Options) (err error) {
	// ...

	return nil
}
