package conv

import (
	"testing"

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
