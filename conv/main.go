package conv

import "github.com/shopspring/decimal"

const (
	places int32 = 6
)

func FloatToDecimal(input float64) decimal.Decimal {
	return decimal.NewFromFloat(input).Truncate(places)
}

func StringToDecimal(input string) (decimal.Decimal, error) {
	dec, err := decimal.NewFromString(input)
	if err != nil {
		return decimal.New(0, 0), err
	}
	return dec.Truncate(places), nil
}

func DecimalToFloat(input decimal.Decimal) float64 {
	result, _ := input.Truncate(places).Float64()
	return result
}

func DecimalToFloatPrecision(input decimal.Decimal, precision int32) float64 {
	result, _ := input.Truncate(precision).Float64()
	return result
}

func Truncate(input decimal.Decimal) decimal.Decimal {
	return input.Truncate(places)
}

func TruncatePrecision(input decimal.Decimal, precision int32) decimal.Decimal {
	return input.Truncate(precision)
}
