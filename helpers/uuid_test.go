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
	assert.False(t, VerifyUUIDv4String(""))
	assert.False(t, VerifyUUIDv4String("empty"))
	assert.False(t, VerifyUUIDv4String(uuid.NewV1().String()))

	assert.True(t, VerifyUUIDv4String("b17c7f4b981e43658679d16d5837a7eb"))
	assert.True(t, VerifyUUIDv4String("b17c7f4b-981e-4365-8679-d16d5837a7eb"))
	assert.False(t, VerifyUUIDv4String("b17c7f4b-981e-6365-8679-d16d5837a7eb"))
	assert.False(t, VerifyUUIDv4String("k17c7f4b-981e-4365-8679-d16d5837a7eb"))

	assert.True(t, VerifyUUIDv4String(GenerateUUIDv4AsString()))
	assert.True(t, VerifyUUIDv4String("{"+GenerateUUIDv4AsString()+"}"))
	assert.False(t, VerifyUUIDv4String("["+GenerateUUIDv4AsString()+"]"))
}

func TestUUIDs(t *testing.T) {
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
				in[i] = GenerateUUIDv4AsString()
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
