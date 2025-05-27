package background

import (
	"backend/internal/models"
	"errors"
	"fmt"
	"log"
	"github.com/wneessen/go-mail"
)

const (
	SMTPServer = "smtp.gmail.com"
	SMTPPort   = 587
)

func SendEmail(notificationConfig models.NotificationConfig, senderConfig models.NotificationSender, notificationRule models.NotificationRule, notificationMessage string) error {
	message := mail.NewMsg()
	if err := message.From(senderConfig.Email); err != nil {
		log.Printf("Failed to set 'From' address: %s", err)
		return errors.New("failed to set 'From' address")
	}
	if err := message.To(notificationConfig.Email); err != nil {
		log.Printf("Failed to set 'To' address: %s", err)
		return errors.New("failed to set 'To' address")
	}

	message.Subject(fmt.Sprintf("Alert - %s", notificationRule.Name))
	message.SetBodyString(mail.TypeTextPlain, notificationMessage)
	client, err := mail.NewClient(
		SMTPServer,
		mail.WithPort(SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(senderConfig.Email),
		mail.WithPassword(senderConfig.AppPassword),
	)
	if err != nil {
		log.Printf("Failed to create mail client: %s", err)
		return errors.New("failed to create mail client")
	}
	if err := client.DialAndSend(message); err != nil {
		log.Printf("Failed to send mail: %s", err)
		return errors.New("failed to send mail")
	}

	return nil
}
