package smtp

import (
	"fmt"
	"learn/settings"
	"net/smtp"
	"strings"
)

func SendMail(to []string, code string, template string) {

	from := settings.Settings.Smtp.From
	host := settings.Settings.Smtp.Host
	password := settings.Settings.Smtp.Password
	port := settings.Settings.Smtp.Port
	auth := smtp.PlainAuth("", from, password, host)

	msg := buildEmailMessage(from, to, from, code)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, from, to, []byte(msg))
	if err != nil {
		panic(err)
	}
}

// buildEmailMessage constructs the email message with headers and body
func buildEmailMessage(from string, to []string, subject string, body string) string {
	// Create the headers
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", strings.Join(to, ", ")),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"UTF-8\"",
	}

	// Join headers with CRLF and add a blank line before the body
	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
