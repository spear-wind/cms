package events

import (
	"strings"
	"testing"
)

func TestNewUserRegistrationEvent(t *testing.T) {
	toEmail := "test@spearwind.io"
	verificationCode := "ABC123"
	event := NewUserRegistrationEvent(toEmail, verificationCode)

	if &event.When == nil {
		t.Errorf("event.When is nil")
	}

	if &event.Message == nil {
		t.Errorf("event.Message is nil")
	}

	if event.Message.To != toEmail {
		t.Errorf("event.Message.To did not equal %s", toEmail)
	}

	if event.Message.Subject != "SpearWind.io - New Account Verification" {
		t.Errorf("User registration event subject was not the expected value: %s", event.Message.Subject)
	}

	if event.Message.From != "no-reply@spearwind.io" {
		t.Errorf("User registration event from email address was not the expected value: %s", event.Message.From)
	}

	emailBody, err := event.Message.Body.String()

	if err != nil {
		t.Errorf("Failed to load user registration event email body: %v", err)
	}

	if strings.Contains(emailBody, "http://spearwind.io/user/verify/"+verificationCode) != true {
		t.Errorf("User registration event email body did not contain the correct user verification link")
	}
}
