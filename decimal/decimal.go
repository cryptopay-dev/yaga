package decimal

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type Decimal = decimal.Decimal

var Zero = decimal.Zero

func NewFromFloat(val float64) Decimal {
	return decimal.NewFromFloat(val)
}

func NewFromString(val string) (Decimal, error) {
	return decimal.NewFromString(val)
}

func New(value int64, exp int32) Decimal {
	return decimal.New(value, exp)
}

func NewFromBigInt(value *big.Int, exp int32) Decimal {
	return decimal.NewFromBigInt(value, exp)
}
