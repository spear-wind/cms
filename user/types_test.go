package user

import "testing"

func TestValidateWithEmptyRequiredFieldsFailsWithErrors(t *testing.T) {
	user := newUser(-1, "", "", "")

	result := user.validate()

	if len(result.Errors) != 4 {
		t.Errorf("Expected four errors, but there were %v errors: %v", len(result.Errors), result.Errors)
	}
}

func TestValidateWithAllRequiredFieldsIsValid(t *testing.T) {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@tld.com",
		Password:  "p@$$w0rd",
	}

	result := user.validate()

	if result.HasErrors() != false {
		t.Fatal("Expected validation to pass with all required fields, but it did not")
	}
}

func TestRegister(t *testing.T) {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@tld.com",
		Password:  "p@$$w0rd",
	}

	result, err := user.register()

	if result.HasErrors() != false {
		t.Fatal("Expected validation to pass with all required fields, but it did not")
	}

	if err != nil {
		t.Errorf("user.register() returned with an unexpected error: %v", err)
	}

	if user.Password != "" {
		t.Errorf("user.Password should be blank")
	}

	if user.VerificationCode == "" {
		t.Error("user.VerificationCode should not be blank")
	}

	if user.hash == "" {
		t.Error("user.hash should not be blank")
	}

	if user.Verified {
		t.Error("user.Verified should be false")
	}
}

func TestVerifyUnverifiedUser(t *testing.T) {
	verificationCode := "ABC123"

	user := User{
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john@tld.com",
		Password:         "",
		VerificationCode: verificationCode,
		Verified:         false,
		hash:             "hashed-password",
	}

	err := user.verify(verificationCode)

	if err != nil {
		t.Errorf("user.verify returned an unexpected error: %v", err)
	}

	if user.VerificationCode != "" {
		t.Errorf("user.VerificationCode should be blank after a successful verification")
	}

	if user.Verified != true {
		t.Errorf("user.Verified should be true after a successful verification")
	}
}

func TestVerifyVerifiedUserReturnsError(t *testing.T) {
	verificationCode := "ABC123"

	user := User{
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john@tld.com",
		Password:         "",
		VerificationCode: verificationCode,
		Verified:         true,
		hash:             "hashed-password",
	}

	err := user.verify(verificationCode)

	if err == nil {
		t.Errorf("Calling verify on a user with .Verified = true should return an error")
	}
}

func TestVerifyWithInvalidVerificationCodeReturnsError(t *testing.T) {
	user := User{
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john@tld.com",
		Password:         "",
		VerificationCode: "ABC123",
		Verified:         false,
		hash:             "hashed-password",
	}

	err := user.verify("321CBA")

	if err == nil {
		t.Errorf("Unmatching verification codes should cause an error")
	}

	if user.VerificationCode == "" {
		t.Errorf("user.VerificationCode should not be blank after a failed verification")
	}

	if user.Verified != false {
		t.Errorf("user.Verified should be false after a failed verification")
	}
}
