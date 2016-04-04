package events

import (
	"fmt"
	"os"
	"time"

	"github.com/dave-malone/email"
)

type EmailEvent struct {
	When    time.Time
	Message email.Message
}

type emailEventSubscriber struct {
	sender email.Sender
}

func (s emailEventSubscriber) Receive(e interface{}) {
	if emailEvent, ok := e.(EmailEvent); ok {
		s.sender.Send(&emailEvent.Message)
	}
}

func NewEmailEventSubscriber(sender email.Sender) EventSubscriber {
	return emailEventSubscriber{sender: sender}
}

func NewUserRegistrationEvent(emailAddress string, verificationCode string) *EmailEvent {
	emailTemplateDir := os.Getenv("EMAIL_TEMPLATE_DIR")
	if emailTemplateDir == "" {
		emailTemplateDir = "../email-templates"
		fmt.Printf("Using %s as the email template directory. Please set env var EMAIL_TEMPLATE_DIR to override this setting", emailTemplateDir)
	}

	messageBody := email.NewFileBasedHTMLTemplateMessageBody(emailTemplateDir+"/user-registration.tpl", verificationCode)

	emailMessage := email.NewMessage(
		"no-reply@spearwind.io",
		emailAddress,
		"SpearWind.io - New Account Verification",
		messageBody,
	)

	emailEvent := &EmailEvent{
		When:    time.Now(),
		Message: *emailMessage,
	}

	return emailEvent
}
