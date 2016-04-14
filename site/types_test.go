package site

import (
	"encoding/json"
	"testing"

	"github.com/spear-wind/cms/user"
)

func TestValidateWithEmptyRequiredFieldsFailsWithErrors(t *testing.T) {
	site := NewSite("", "", nil)

	result := site.validate()

	if len(result.Errors) != 3 {
		t.Errorf("Expected exactly three errors, but there were %d errors: %v", len(result.Errors), result.Errors)
	}
}

func TestValidateHappyPath(t *testing.T) {
	user := &user.User{
		FirstName: "Spearwind",
		LastName:  "Creator",
		Email:     "creator@spearwind.io",
		Password:  "p@$$w0rd",
	}

	site := NewSite("Spearwind", "Spearwind.io", user)

	result := site.validate()

	if result.HasErrors() != false {
		t.Error("Expected validation to pass with all required fields, but it did not")
	}
}

func TestMarshallJSONHappyPath(t *testing.T) {
	b := []byte(`{"name":"SpearWind","domain_name":"spearwind.io"}`)
	var site Site

	if err := json.Unmarshal(b, &site); err != nil {
		t.Errorf("Unexpected error unmarshalling json to Site: %v", err)
	}

	if len(site.Name) == 0 {
		t.Error("json name field failed to unmarshall to Site.Name")
	}

	if len(site.DomainName) == 0 {
		t.Error("json domain_name field failed to unmarshall to Site.DomainName")
	}
}

func TestMarshallBadJSONResultsInSiteWithUninitializedFields(t *testing.T) {
	b := []byte(`{"bad":"bad json!"}`)
	var site Site

	if err := json.Unmarshal(b, &site); err != nil {
		t.Errorf("Unexpected error unmarshalling json to Site: %v", err)
	}

	if len(site.Name) != 0 {
		t.Error("bad json should result in uninitialized Site.Name field")
	}

	if len(site.DomainName) != 0 {
		t.Error("bad json should result in uninitialized Site.DomainName field")
	}
}
