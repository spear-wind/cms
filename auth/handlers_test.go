package auth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func TestLoginHandlerResposnseToInvalidData(t *testing.T) {
	client := &http.Client{}
	userRepository := user.NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(loginHandler(formatter, userRepository)))

	form := url.Values{}
	form.Add("foo", "asdf")
	form.Add("bar", "1234")

	req, err := http.NewRequest("POST", server.URL, strings.NewReader(form.Encode()))
	if err != nil {
		t.Errorf("Error in creating POST request with loginHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to loginHandler: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending bad form data should result in a BadRequest response code from the server")
	}
}

func TestLoginHandlerHappyPath(t *testing.T) {
	client := &http.Client{}
	userRepository := user.NewInMemoryRepository()

	user := user.User{
		FirstName: "Adam",
		LastName:  "Spearwind",
		Email:     "test@spearwind.io",
		Password:  "abc123",
	}

	validationResult, err := user.Register()
	if validationResult.HasErrors() != false {
		t.Fatal("Expected user validation to pass with all required fields, but it did not")
	}

	if err != nil {
		t.Fatalf("user.Register() returned with an unexpected error: %v", err)
	}

	userRepository.Add(user)

	server := httptest.NewServer(http.HandlerFunc(loginHandler(formatter, userRepository)))

	form := url.Values{}
	form.Add("email", "test@spearwind.io")
	form.Add("password", "abc123")

	req, err := http.NewRequest("POST", server.URL, strings.NewReader(form.Encode()))
	if err != nil {
		t.Errorf("Error in creating POST request with loginHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to loginHandler: %v", err)
	}
	defer res.Body.Close()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Successful login should result in http.StatusOK; \nstatus code: %v\nresponse body: %s", res.StatusCode, payload)
	}
}
