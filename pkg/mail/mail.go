package mail

import (
	gomail "gopkg.in/gomail.v2"
)

type Mailer interface {
	Send(to string, subject string, bodyType string, body string) error
}

type Mail struct {
	Port int
	Host string

	To []string

	From     string
	Username string
	Password string
}

func NewMail(Host string, Port int, From string, To []string, Username string, Password string) *Mail {
	return &Mail{
		Port,
		Host,
		To,
		From,
		Username,
		Password,
	}
}

func (m *Mail) Send(to string, subject string, bodyType string, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody(bodyType, body)
	n := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)

	return n.DialAndSend(msg)
}
