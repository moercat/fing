package email

import (
	"fing/pkg/config"
	"github.com/go-gomail/gomail"
)

func SendMail(toAddress, name, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", config.Config.Email.Email, config.Config.Email.Name)
	msg.SetHeader("To", msg.FormatAddress(toAddress, name))
	msg.SetHeader("Subject", subject) // 主题
	msg.SetBody("text/html", body)
	dialer := gomail.NewDialer(config.Config.Email.Host, 465, config.Config.Email.Email, config.Config.Email.Password)
	return dialer.DialAndSend(msg)
}
