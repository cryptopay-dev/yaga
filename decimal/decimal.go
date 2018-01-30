package decimal

import (
	"math/big"

	"github.com/shopspring/decimal"
)

// Decimal type alias for basic decimal.Decimal
type Decimal = decimal.Decimal

// Zero type alias for basic decimal.Zero
var Zero = decimal.Zero

// NewFromFloat return decimal from float number
func NewFromFloat(val float64) Decimal {
	return decimal.NewFromFloat(val)
}

// NewFromString return decimal and error from number in string
func NewFromString(val string) (Decimal, error) {
	return decimal.NewFromString(val)
}

// New create new decimal with value and exp
func New(value int64, exp int32) Decimal {
	return decimal.New(value, exp)
}

// NewFromBigInt create decimal number from big int
func NewFromBigInt(value *big.Int, exp int32) Decimal {
	return decimal.NewFromBigInt(value, exp)
}
