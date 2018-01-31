package conv_test

import (
	"fmt"

	"github.com/cryptopay-dev/yaga/conv"
	"github.com/cryptopay-dev/yaga/decimal"
)

func ExampleDecimalToFloat() {
	f := conv.DecimalToFloat(decimal.NewFromFloat(0.004))
	fmt.Println(f)

	// Output:
	// 0.004
}

func ExampleStringToDecimal() {
	d, err := conv.StringToDecimal("0.123")
	if err != nil {
		panic(err)
	}
	fmt.Println(d)

	// Output:
	// 0.123
}

func ExampleDecimalToFloatPrecision() {
	f := conv.DecimalToFloatPrecision(decimal.NewFromFloat(0.987), int32(2))
	fmt.Println(f)

	// Output:
	// 0.98
}
