package conv

import (
	"testing"

	"github.com/cryptopay-dev/yaga/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFloatToDecimal(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{10.000000110, 10.0},
		{0.0000000000001, 0},
		{15.001, 15.001},
		{0.0006, 0.0006},
		{0.0005678, 0.000567},
	}

	for _, test := range tests {
		v := FloatToDecimal(test.value)
		res := DecimalToFloat(v)
		assert.Equal(t, test.expected, res)
	}
}

func TestStringToDecimal(t *testing.T) {
	s, err := StringToDecimal("0.002")
	assert.NoError(t, err)
	f, _ := s.Float64()
	assert.Equal(t, 0.002, f)
}

func TestWrongStringToDecimal(t *testing.T) {
	_, err := StringToDecimal("wrong string")
	assert.Error(t, err)
}

func TestDecimalToFloatPrecision(t *testing.T) {
	f := DecimalToFloatPrecision(decimal.NewFromFloat(0.59), int32(1))
	assert.Equal(t, 0.5, f)
}

func TestTruncatePrecision(t *testing.T) {
	d := TruncatePrecision(decimal.NewFromFloat(0.5059), int32(3))
	assert.Equal(t, decimal.NewFromFloat(0.505), d)
}
