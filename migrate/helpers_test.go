package migrate

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractAttributes(t *testing.T) {
	var tests = []struct {
		actual  string
		version int64
		name    string
		mType   string
		valid   bool
	}{
		{
			actual:  "1512662684_auth.up.sql",
			version: 1512662684,
			name:    "auth",
			mType:   "up",
			valid:   true,
		},

		{
			actual:  "1512662684_auth.down.sql",
			version: 1512662684,
			name:    "auth",
			mType:   "down",
			valid:   true,
		},

		{
			actual:  "1512662684_something_like_name.down.sql",
			version: 1512662684,
			name:    "something_like_name",
			mType:   "down",
			valid:   true,
		},

		{
			actual: "1512662684something_like_name.down.sql",
			valid:  false,
		},

		{
			actual: "something_like_name.down.sql",
			valid:  false,
		},

		{
			actual: "1512662684_something_like_name.sql",
			valid:  false,
		},

		{
			actual: "-1512662684_something_like_name.down.sql",
			valid:  false,
		},

		{
			actual: "",
			valid:  false,
		},
	}

	for _, item := range tests {
		version, name, mType, err := extractAttributes(item.actual)

		if item.valid {
			assert.NoError(t, err)
			assert.Equal(t, item.version, version)
			assert.Equal(t, item.name, name)
			assert.Equal(t, item.mType, mType)
		} else {
			assert.Errorf(t, err, "It must fail: %s", item.actual)
		}
	}
}

func TestExtractMigrations(t *testing.T) {
	var (
		err   error
		path  = "./fixtures/good"
		files []os.FileInfo
	)

	if files, err = findMigrations(path); !assert.NoError(t, err) {
		t.FailNow()
	}

	items, errMigrate := extractMigrations(defaultLogger, path, files)

	if !assert.NoError(t, errMigrate) {
		t.FailNow()
	}

	// up/down is one migration:
	assert.Equal(t, len(files)/2, len(items))
}

func TestCreateMigrations(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		dir, err := ioutil.TempDir("", time.Now().Format(time.RFC3339))
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		if err = CreateMigration(dir, "something"); !assert.NoError(t, err) {
			t.FailNow()
		}

		files, err := findMigrations(dir)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		if !assert.True(t, 2 == len(files)) {
			t.FailNow()
		}

		items, err := extractMigrations(defaultLogger, dir, files)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		if !assert.True(t, 1 == len(items)) {
			t.FailNow()
		}

		if !assert.True(t, "something" == items[0].Name) {
			t.FailNow()
		}
	})

	t.Run("Bad", func(t *testing.T) {
		if err := CreateMigration("/path/to/not/existing/folder", "something"); !assert.Error(t, err) {
			t.FailNow()
		}
	})
}
