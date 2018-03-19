package mail

import (
	"github.com/mattbaird/gochimp"
)

type Events struct {
	mail       *mailService
	subject    string
	recipients []gochimp.Recipient
}

func (e *Events) Send(message string) {
	e.mail.msgCh <- gochimp.Message{
		Text:    message,
		Subject: e.subject,
		To:      e.recipients,
	}
}
