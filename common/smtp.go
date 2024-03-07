package common

import (
	"bytes"
	"fmt"
	"net/smtp"
)

type ISmtpService interface {
	SendEmail(subject string, to []string, msg []byte) error
}

type SmtpService struct {
	username string
	password string
	smtpHost string
	smtpPort int
	from     string
}

func NewSmtpService(from, username, password, smtpHost string, smtpPort int) *SmtpService {
	return &SmtpService{
		username: username,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
	}
}

func (s *SmtpService) SendEmail(subject string, to []string, msg []byte) error {
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))
	body.Write(msg)
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	return smtp.SendMail(fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort), auth, s.from, to, body.Bytes())
}
