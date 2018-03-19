package mail

import (
	"time"
)

type event struct {
	lastSendOk  time.Time
	lastSendErr time.Time
}
