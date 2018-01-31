package main

import (
	"fmt"

	"github.com/cryptopay-dev/yaga/config"
)

// Config structure
type Config struct {
	A int `yaml:"a" validate:"required,gt=0,eq=1"`
	B int `yaml:"b" validate:"required,gt=0,eq=2"`
	C int `yaml:"c" validate:"required,gt=0,eq=3"`
	D int `yaml:"d" validate:"required,gt=0,eq=4"`
}

func main() {
	var conf Config

	if err := config.Load("config.fixture.yaml", &conf); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", conf)

	// Output:
	// main.Config{A:1, B:2, C:3, D:4}
}
