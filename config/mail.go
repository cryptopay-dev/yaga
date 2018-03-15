package config

import (
	"time"

	"github.com/cryptopay-dev/yaga/mail"
)

// Mail default configuration
type Mail struct {
	APIKey            string        `yaml:"api_key" validate:"required"`
	Recipients        []string      `yaml:"recipients" validate:"required,min=1"`
	FromEmail         string        `yaml:"from_email" validate:"required,email"`
	FromName          string        `yaml:"from_name" validate:"required"`
	SendUniqTimeout   time.Duration `yaml:"send_uniq_timeout" validate:"gte=1"`
	RetryErrorTimeout time.Duration `yaml:"retry_error_timeout" validate:"gte=1"`
}

// Connect to Mail
func (g Mail) Connect(log mail.Logger) (mail.Mailer, error) {
	return mail.New(mail.Options{
		APIKey:            g.APIKey,
		Logger:            log,
		Recipients:        g.Recipients,
		FromEmail:         g.FromEmail,
		FromName:          g.FromName,
		SendUniqTimeout:   g.SendUniqTimeout,
		RetryErrorTimeout: g.RetryErrorTimeout,
	})
}
