package helpers

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// GenerateUUIDv4AsString returns random generated canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func GenerateUUIDv4AsString() string {
	return uuid.NewV4().String()
}

// VerifyUUIDv4String verification UUID as string for conformity with version 4
func VerifyUUIDv4String(src string) bool {
	if len(src) < uuid.Size*2 {
		return false
	}

	if uid, err := uuid.FromString(src); err != nil {
		return false
	} else if uid.Version() != 4 {
		return false
	}

	return true
}

// ValidateUUIDs validation UUIDs for conformity with version 4
func ValidateUUIDs(uuids []string) (err error) {
	var uid uuid.UUID

	if len(uuids) == 0 {
		return
	}

	for _, id := range uuids {
		if uid, err = uuid.FromString(id); err != nil {
			return fmt.Errorf("uuid: validation error(`%s`): %v", id, err)
		}

		if uid.Version() != 4 {
			return fmt.Errorf("uuid: invalid version number (`%s`): %d", id, uid.Version())
		}
	}

	return
}
