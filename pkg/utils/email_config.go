package utils

import "fmt"

type EmailConfig struct {
	SMTPHost string
	Port     string
	From     string
	To       string
	Subject  string
	Body     string
}

func (config *EmailConfig) ValidateEmailConfig() error {
	if config.SMTPHost == "" || config.Port == "" || config.From == "" || config.To == "" {
		return fmt.Errorf("invalid email configuration: all fields are required")
	}
	return nil
}
