package mailer

import (
	"embed"
	"github.com/go-mail/mail/v2"
	"time"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	sender string
	dialer *mail.Dialer
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		sender: sender,
		dialer: dialer,
	}
}
