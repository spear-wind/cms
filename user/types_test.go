package user

import "testing"

func TestValidateWithEmptyRequiredFieldsFailsWithErrors(t *testing.T) {
	user := newUser(-1, "", "", "")

	result := user.validate()

	if len(result.Errors) != 3 {
		t.Fatalf("Expected three errors, but there were %v errors: %v", len(result.Errors), result.Errors)
	}
}

func TestValidateWithAllRequiredFieldsIsValid(t *testing.T) {
	user := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@tld.com",
	}

	result := user.validate()

	if result.HasErrors() != false {
		t.Fatal("Expected validation to pass with all required fields, but it did not")
	}
}
