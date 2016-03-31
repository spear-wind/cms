package validator

import "testing"

func TestValidateEmailHappyPath(t *testing.T) {
	email := "yo@yo.com"
	valid := ValidateEmail(email)

	if valid != true {
		t.Fatalf("%s was not valid", email)
	}
}
func TestValidateEmailWithMixedCaseAlphaNumericPasses(t *testing.T) {
	email := "yaYO123@hereIam3.com"
	valid := ValidateEmail(email)

	if valid != true {
		t.Fatalf("%s was not valid", email)
	}
}

func TestValidateInvalidEmailMissingDomainExtensionFails(t *testing.T) {
	email := "yo@yo"
	valid := ValidateEmail(email)

	if valid {
		t.Fatalf("%s should not pass validation", email)
	}
}

func TestValidateLongTldFails(t *testing.T) {
	email := "yo@yo.comcomcom"
	valid := ValidateEmail(email)

	if valid {
		t.Fatalf("%s should not pass validation", email)
	}
}
