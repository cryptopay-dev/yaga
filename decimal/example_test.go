package decimal_test

import (
	"fmt"
	"math/big"

	"github.com/cryptopay-dev/yaga/decimal"
)

func ExampleNewFromFloat() {
	d := decimal.NewFromFloat(0.5)
	fmt.Println(d)

	// Output:
	// 0.5
}

func ExampleNewFromString() {
	d, err := decimal.NewFromString("0.234")
	if err != nil {
		panic(err)
	}
	fmt.Println(d)

	// Output:
	// 0.234
}

func ExampleNewFromBigInt() {
	var bigInt big.Int
	bigInt.SetInt64(777)
	d := decimal.NewFromBigInt(&bigInt, 0)
	fmt.Println(d.String())

	// Output:
	// 777
}
