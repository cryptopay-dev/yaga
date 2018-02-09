package helpers

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDv4(t *testing.T) {
	assert.False(t, VerifyUUIDv4String(""))
	assert.False(t, VerifyUUIDv4String("empty"))
	assert.False(t, VerifyUUIDv4String(uuid.NewV1().String()))

	assert.True(t, VerifyUUIDv4String("b17c7f4b981e43658679d16d5837a7eb"))
	assert.True(t, VerifyUUIDv4String("b17c7f4b-981e-4365-8679-d16d5837a7eb"))
	assert.False(t, VerifyUUIDv4String("k17c7f4b-981e-4365-8679-d16d5837a7eb"))

	assert.True(t, VerifyUUIDv4String(GenerateUUIDv4AsString()))
	assert.True(t, VerifyUUIDv4String("{"+GenerateUUIDv4AsString()+"}"))
	assert.False(t, VerifyUUIDv4String("["+GenerateUUIDv4AsString()+"]"))
}
