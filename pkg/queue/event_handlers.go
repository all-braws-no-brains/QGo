package queue

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/all-braws-no-brains/QGo/pkg/utils"
)

// EmailEventHandler handles email-related events
func EmailEventHandler(ctx context.Context, event *Event, config *utils.EmailConfig) error {
	if err := config.ValidateEmailConfig(); err != nil {
		return fmt.Errorf("invalid email configuration: %v", err)
	}

	// Type assertion to ensure the payload is a string (email body content in this case)
	payload, ok := event.Payload.(string)
	if !ok {
		return fmt.Errorf("invalid payload type: expected string, got %T", event.Payload)
	}
	emailContent := event.Payload.(string)

	// Use the config to get the SMTP settings.
	smtpHost := config.SMTPHost
	port := config.Port
	from := config.From
	to := config.To
	subject := "Subject: " + config.Subject
	body := "Content-Type: text/plain; charset=UTF-8\r\n\r\n" + emailContent

	// Setup SMTP server
	auth := smtp.PlainAuth("", from, "your-password", smtpHost)

	// Build the email message
	message := []byte(subject + "\r\n" + body)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+port, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	// Simulate email sent log
	fmt.Printf("Email sent successfully for event %s: %s\n", event.ID, payload)
	return nil
}
