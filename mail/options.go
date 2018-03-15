package mail

import (
	"time"
)

type Options struct {
	APIKey            string
	Logger            Logger
	Recipients        []string
	FromEmail         string
	FromName          string
	SendUniqTimeout   time.Duration
	RetryErrorTimeout time.Duration
}
