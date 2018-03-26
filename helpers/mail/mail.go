package mail

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/mail"
)

// Connect to Mail
func Connect(key string, log mail.Logger) (mail.Mailer, error) {
	return mail.New(mail.Options{
		APIKey:            config.GetString(key + ".api_key"),
		Logger:            log,
		Recipients:        config.GetStringSlice(key + ".recipients"),
		FromEmail:         config.GetString(key + ".from_email"),
		FromName:          config.GetString(key + ".from_name"),
		SendUniqTimeout:   config.GetDuration(key + ".send_timeout"),
		RetryErrorTimeout: config.GetDuration(key + ".retry_timeout"),
	})
}
