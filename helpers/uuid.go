package helpers

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// NewUUIDv4 returns random generated canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func NewUUIDv4() string {
	return uuid.NewV4().String()
}

// ValidateUUIDv4 verification UUID as string for conformity with version 4
func ValidateUUIDv4(src string) error {
	if len(src) < uuid.Size*2 {
		return fmt.Errorf("uuid: incorrect UUID length: %s", src)
	}

	if uid, err := uuid.FromString(src); err != nil {
		return err
	} else if uid.Version() != 4 {
		return fmt.Errorf("uuid: invalid version number, must be '4', actual '%d'", uid.Version())
	}

	return nil
}

// ValidateUUIDs validation UUIDs for conformity with version 4
func ValidateUUIDs(uuids []string) (err error) {
	if len(uuids) == 0 {
		return
	}

	for _, id := range uuids {
		if err = ValidateUUIDv4(id); err != nil {
			return
		}
	}

	return
}
