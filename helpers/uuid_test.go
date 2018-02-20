package helpers_test

import (
	"testing"

	. "github.com/cryptopay-dev/yaga/helpers"
	"github.com/stretchr/testify/assert"
)

var uuids = []struct {
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
		// UUIDv1
		uuid:   "87adb170-1568-11e8-b642-0ed5f89f718b",
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

var okUUIDs = []string{
	"de82cc12-2d06-4fb2-8446-83deb13e4f46",
	"7db6a668-3d97-440e-99d1-68f902d6e6bf",
	"876b07bb-9008-47f5-857f-e0170923f6c7",
	"d0732409-cf5f-43e6-8673-d3ff35ceb76f",
	"b9e7117e-4fb4-4c16-9b73-1844359b2944",
	"b610263c-a361-4c15-b17f-2152ee115f0a",
	"fbdf5285-5f52-4db5-b2ec-6c3c8709711a",
	"afd62dca-6269-46b5-a1d0-0d518b4a2ebd",
	"b32d1065-6bf6-464c-a011-8b32e9961924",
	"32c2c280-df3f-4d98-8030-f59859e87835",
}

var positionsUUID = []int{0, len(okUUIDs) / 2, len(okUUIDs) - 1}

func TestNewUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		assert.NoError(t, ValidateUUID(NewUUID()))
	}

	assert.NoError(t, ValidateUUID("{"+NewUUID()+"}"))
	assert.Error(t, ValidateUUID("["+NewUUID()+"]"))
}

func TestValidateUUID(t *testing.T) {
	assert.Error(t, ValidateUUID(""))
	assert.Error(t, ValidateUUID("empty"))
	assert.Error(t, ValidateUUID("b610263c-a361-1c15-b17f-2152ee115f0a")) // UUIDv1

	assert.NoError(t, ValidateUUID("b17c7f4b981e43658679d16d5837a7eb"))
	assert.NoError(t, ValidateUUID("b17c7f4b-981e-4365-8679-d16d5837a7eb"))
	assert.Error(t, ValidateUUID("b17c7f4b-981e-6365-8679-d16d5837a7eb")) // not UUIDv4
	assert.Error(t, ValidateUUID("k17c7f4b-981e-4365-8679-d16d5837a7eb")) // not valid symbols

	for _, id := range okUUIDs {
		assert.NoError(t, ValidateUUID(id))
	}

	for _, item := range uuids {
		if item.result {
			assert.NoError(t, ValidateUUID(item.uuid))
		} else {
			assert.Error(t, ValidateUUID(item.uuid))
		}
	}
}

func TestValidateUUIDs(t *testing.T) {
	var (
		err error

		size = len(okUUIDs)
	)

	for _, k := range positionsUUID {
		// k - insertion position
		for _, item := range uuids {
			in := make([]string, size)
			for i := 0; i < size; i++ {
				if i == k {
					in[i] = item.uuid
				} else {
					in[i] = okUUIDs[i]
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
}
