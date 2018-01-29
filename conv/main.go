package conv

import "github.com/cryptopay-dev/yaga/decimal"

const (
	places int32 = 6
)

// FloatToDecimal convert float to decimal
func FloatToDecimal(input float64) decimal.Decimal {
	return decimal.NewFromFloat(input).Truncate(places)
}

// StringToDecimal convert number in string to decimal
func StringToDecimal(input string) (decimal.Decimal, error) {
	dec, err := decimal.NewFromString(input)
	if err != nil {
		return decimal.New(0, 0), err
	}
	return dec.Truncate(places), nil
}

// DecimalToFloat convert decimal to float
func DecimalToFloat(input decimal.Decimal) float64 {
	result, _ := input.Truncate(places).Float64()
	return result
}

// DecimalToFloatPrecision convert decimal to float with precision
func DecimalToFloatPrecision(input decimal.Decimal, precision int32) float64 {
	result, _ := input.Truncate(precision).Float64()
	return result
}

// Truncate truncate decimal number with default precision
func Truncate(input decimal.Decimal) decimal.Decimal {
	return input.Truncate(places)
}

// TruncatePrecision truncate decimal number with param precision
func TruncatePrecision(input decimal.Decimal, precision int32) decimal.Decimal {
	return input.Truncate(precision)
}
