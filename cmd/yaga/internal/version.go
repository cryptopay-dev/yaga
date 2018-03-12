package internal

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// VersionFile to parse
var VersionFile = path.Join(os.Getenv("GOPATH"), "src", "github.com/cryptopay-dev/yaga/cmd/yaga/VERSION")

func parseVersion(f io.Reader) (string, string, error) {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", "", err
	}

	version := strings.Split(string(data), "\n")

	if len(version) < 2 {
		return "", "", errors.New("bad version file")
	}

	return version[0], version[1], nil
}

// Version method parse file
func Version() (string, string, error) {
	file, err := os.Open(VersionFile)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	return parseVersion(file)
}
