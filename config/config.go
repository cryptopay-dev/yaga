package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

var ErrUnknownSourceType = errors.New("unknown type")

func Load(src, config interface{}) error {
	switch t := src.(type) {
	case io.Reader:
		return readAndValidate(t.(io.Reader), config)
	case string:
		file, err := os.Open(t)
		if err != nil {
			return err
		}
		return readAndValidate(file, config)
	default:
		return ErrUnknownSourceType
	}
}

func readAndValidate(reader io.Reader, config interface{}) error {
	var (
		err error
		buf []byte
	)

	if buf, err = ioutil.ReadAll(reader); err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, config); err != nil {
		return err
	}

	return validator.New().Struct(config)
}
