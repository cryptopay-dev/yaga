package mail

import (
	"time"

	"github.com/mattbaird/gochimp"
)

type Mailer interface {
	Events(subject string, recipients []string) *Events
}

type mailService struct {
	msgCh      chan gochimp.Message
	events     map[string]*event
	options    Options
	recipients []gochimp.Recipient
	api        *gochimp.MandrillAPI
}

func New(opts Options) (Mailer, error) {
	var err error
	m := &mailService{
		events:     make(map[string]*event),
		msgCh:      make(chan gochimp.Message, 1),
		options:    opts,
		recipients: formatRecipients(opts.Recipients),
	}

	m.api, err = gochimp.NewMandrill(opts.APIKey)
	if err != nil {
		return nil, err
	}
	go m.worker()

	return m, nil
}

func formatRecipients(emails []string) []gochimp.Recipient {
	recipients := make([]gochimp.Recipient, 0, len(emails))
	for _, email := range emails {
		recipients = append(recipients, gochimp.Recipient{Email: email})
	}
	return recipients
}

func (m *mailService) Events(subject string, recipients []string) *Events {
	return &Events{
		mail:       m,
		subject:    subject,
		recipients: formatRecipients(recipients),
	}
}

func (m *mailService) worker() {
	var (
		e      *event
		err    error
		found  bool
		now, t time.Time
	)

	for msg := range m.msgCh {
		now = time.Now()
		if e, found = m.events[msg.Text]; found {
			if !e.lastSendOk.IsZero() {
				t = e.lastSendOk.Add(m.options.SendUniqTimeout)
			} else {
				t = e.lastSendErr.Add(m.options.RetryErrorTimeout)
			}
			if now.Before(t) {
				continue
			}
		} else {
			e = new(event)
			m.events[msg.Text] = e
		}

		msg.FromEmail = m.options.FromEmail
		msg.FromName = m.options.FromName
		if len(msg.To) == 0 {
			msg.To = m.recipients
		}

		_, err = m.api.MessageSend(msg, true)
		if err != nil {
			e.lastSendErr = now
			e.lastSendOk = time.Time{}
			if m.options.Logger != nil {
				m.options.Logger.Errorf("Error sending message: %s", msg.Subject)
			}
		} else {
			e.lastSendOk = now
			e.lastSendErr = time.Time{}
		}
	}
}
