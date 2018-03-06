package config

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

type testConfig struct {
	A int `yaml:"a" validate:"required,gt=0,eq=1"`
	B int `yaml:"b" validate:"required,gt=0,eq=2"`
	C int `yaml:"c" validate:"required,gt=0,eq=3"`
	D int `yaml:"d" validate:"required,gt=0,eq=4"`
}

var stringReader = strings.NewReader(`
b: -1
c: -2
d: -3
`)

func TestLoad(t *testing.T) {
	t.Run("should fail on empty", func(t *testing.T) {
		var (
			err  error
			buf  = bytes.NewBuffer(nil)
			conf testConfig
		)

		err = Load(buf, &conf)

		if err != io.EOF {
			t.Fatal(err)
		}
	})

	t.Run("should fail on closed file", func(t *testing.T) {
		var (
			err  error
			file *os.File
			conf testConfig
		)

		if file, err = os.Open("config.fixture.yaml"); err != nil {
			t.Fatal(err)
		}

		file.Close()

		err = Load(file, &conf)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), os.ErrClosed.Error())
		}
	})

	t.Run("should fail on file not exists", func(t *testing.T) {
		var (
			err  error
			conf testConfig
		)

		err = Load("config.unknown.yaml", &conf)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "no such file or directory")
		}
	})

	t.Run("should fail on bad yaml", func(t *testing.T) {
		var (
			err  error
			buf  = bytes.NewBufferString(`a.a.a. "test": "not yaml" \c\c\t}`)
			conf testConfig
		)

		err = Load(buf, &conf)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "yaml: ")
		}
	})

	t.Run("should work with string as path", func(t *testing.T) {
		var (
			path = "config.fixture.yaml"
			conf testConfig
		)

		if err := Load(path, &conf); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testConfig{
			A: 1,
			B: 2,
			C: 3,
			D: 4,
		}, conf)
	})

	t.Run("should work with io.Reader", func(t *testing.T) {
		var (
			err  error
			file *os.File
			conf testConfig
		)

		if file, err = os.Open("config.fixture.yaml"); err != nil {
			t.Fatal(err)
		}

		if err = Load(file, &conf); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testConfig{
			A: 1,
			B: 2,
			C: 3,
			D: 4,
		}, conf)
	})

	t.Run("should fail on unknown source", func(t *testing.T) {
		var conf testConfig
		assert.EqualError(t, Load(nil, &conf), ErrUnknownSourceType.Error())
	})

	t.Run("should validate structure", func(t *testing.T) {
		var (
			err         error
			conf        testConfig
			fieldsCount = reflect.ValueOf(conf).NumField()
		)

		err = Load(stringReader, &conf)

		switch errFields := err.(type) {
		case validator.ValidationErrors:
			assert.Equal(t, fieldsCount, len(errFields))
		default:
			t.Fatal(err)
		}
	})
}
