package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		_, _, err := Version()
		assert.NoError(t, err)
	})

	t.Run("bad", func(t *testing.T) {
		defer func(ver string) {
			VersionFile = ver
		}(VersionFile)

		VersionFile += ".tmp"
		_, _, err := Version()
		assert.Error(t, err)
	})

	t.Run("testcases", func(t *testing.T) {
		t.Run("#1", func(t *testing.T) {
			v, d, err := parseVersion(strings.NewReader("1\n2"))
			assert.NoError(t, err)
			assert.Equal(t, "1", v)
			assert.Equal(t, "2", d)
		})

		t.Run("#2", func(t *testing.T) {
			v, d, err := parseVersion(strings.NewReader("3\n4\n\n"))
			assert.NoError(t, err)
			assert.Equal(t, "3", v)
			assert.Equal(t, "4", d)
		})

		t.Run("#3", func(t *testing.T) {
			_, _, err := parseVersion(strings.NewReader(""))
			assert.Error(t, err)
		})

		t.Run("#4", func(t *testing.T) {
			f, err := os.Open(VersionFile)
			assert.NoError(t, err)
			f.Close()
			_, _, err = parseVersion(f)
			assert.Error(t, err)
		})
	})
}
