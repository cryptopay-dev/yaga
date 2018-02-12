package helpers_test

import (
	"math/rand"
	"testing"

	. "github.com/cryptopay-dev/yaga/helpers"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var uuidsForTest = []struct {
	uuid   string
	result bool
}{
	{
		uuid:   "",
		result: false,
	},
	{
		uuid:   "empty",
		result: false,
	},
	{
		uuid:   uuid.NewV1().String(),
		result: false,
	},
	{
		uuid:   "b17c7f4b981e43658679d16d5837a7eb",
		result: true,
	},
	{
		uuid:   "b17c7f4b-981e-4365-8679-d16d5837a7eb",
		result: true,
	},
	{
		uuid:   "{b17c7f4b-981e-4365-8679-d16d5837a7eb}",
		result: true,
	},
	{
		// not 4 version
		uuid:   "b17c7f4b-981e-7365-8679-d16d5837a7eb",
		result: false,
	},
	{
		uuid:   "b17c7f4b-981e-4365-8x79-d16d5837a7eb",
		result: false,
	},
	{
		uuid:   "k17c7f4b-981e-4365-8679-d16d5837a7eb",
		result: false,
	},
}

func TestUUIDv4(t *testing.T) {
	assert.Error(t, ValidateUUIDv4(""))
	assert.Error(t, ValidateUUIDv4("empty"))
	assert.Error(t, ValidateUUIDv4(uuid.NewV1().String()))

	assert.NoError(t, ValidateUUIDv4("b17c7f4b981e43658679d16d5837a7eb"))
	assert.NoError(t, ValidateUUIDv4("b17c7f4b-981e-4365-8679-d16d5837a7eb"))
	assert.Error(t, ValidateUUIDv4("b17c7f4b-981e-6365-8679-d16d5837a7eb"))
	assert.Error(t, ValidateUUIDv4("k17c7f4b-981e-4365-8679-d16d5837a7eb"))

	assert.NoError(t, ValidateUUIDv4(NewUUIDv4()))
	assert.NoError(t, ValidateUUIDv4("{"+NewUUIDv4()+"}"))
	assert.Error(t, ValidateUUIDv4("["+NewUUIDv4()+"]"))
}

func TestValidateUUIDs(t *testing.T) {
	var (
		err error
		ins int

		size = 5
	)

	for _, item := range uuidsForTest {
		ins = rand.Intn(size)
		in := make([]string, size)
		for i := 0; i < size; i++ {
			if i == ins {
				in[i] = item.uuid
			} else {
				in[i] = NewUUIDv4()
			}
		}

		err = ValidateUUIDs(in)
		if item.result {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
