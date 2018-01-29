package decimal

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatToDecimal(t *testing.T) {
	d := NewFromFloat(0.5)
	assert.Equal(t, "0.5", d.String())
}

func TestFromString(t *testing.T) {
	d, err := NewFromString("0.004")
	if assert.NoError(t, err) {
		f, _ := d.Float64()
		assert.Equal(t, 0.004, f)
	}
}

func TestFromWrongString(t *testing.T) {
	_, err := NewFromString("wrong string")
	assert.Error(t, err)
}

func TestNewFromBigInt(t *testing.T) {
	var bigInt big.Int
	bigInt.SetInt64(777)
	d := NewFromBigInt(&bigInt, 0)
	assert.Equal(t, "777", d.String())
}
