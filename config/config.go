package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

// ErrUnknownSourceType when source type not string-path or io.Reader
var ErrUnknownSourceType = errors.New("unknown type")

// Load config from source (like io.Reader / string path)
// and validate it (https://github.com/go-playground/validator)
func Load(src, config interface{}) error {
	switch t := src.(type) {
	case io.Reader:
		return readAndValidate(t, config)
	case string:
		file, err := os.Open(t)
		if err != nil {
			return err
		}
		defer file.Close()
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
