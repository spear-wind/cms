package facebook

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

type fakeClient struct {
	user *user.User
	err  error
}

func (c *fakeClient) getUser(cmd loginCommand) (*user.User, error) {
	return c.user, c.err
}

func (c *fakeClient) getUserReturns(user *user.User, error error) {
	c.user = user
	c.err = error
}

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func TestLoginHandlerResponseToInvalidData(t *testing.T) {
	client := &http.Client{}
	userRepository := user.NewInMemoryRepository()
	fbClient := new(fakeClient)

	server := httptest.NewServer(http.HandlerFunc(facebookLoginHandler(formatter, userRepository, fbClient)))
	defer server.Close()

	invalidBody := []byte("not even json")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(invalidBody))
	if err != nil {
		t.Errorf("Error in creating POST request with facebookLoginHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to facebookLoginHandler: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending invalid Body should result in a bad request from the server")
	}
}

func TestLoginHandlerResponseToBadJson(t *testing.T) {
	client := &http.Client{}
	userRepository := user.NewInMemoryRepository()
	fbClient := new(fakeClient)

	server := httptest.NewServer(http.HandlerFunc(facebookLoginHandler(formatter, userRepository, fbClient)))
	defer server.Close()

	badJSON := []byte("{\"test\":\"bad json! bad!\"}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(badJSON))
	if err != nil {
		t.Errorf("Error in creating POST request with facebookLoginHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to facebookLoginHandler: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending bad JSON should result in a bad request from the server")
	}
}

func TestLoginHandlerHappyPath(t *testing.T) {
	client := &http.Client{}
	userRepository := user.NewInMemoryRepository()
	fbClient := new(fakeClient)

	fakeUser := &user.User{
		FirstName:  "Test",
		LastName:   "User",
		Email:      "testuser@spearwind.io",
		FacebookID: "987",
	}

	fbClient.getUserReturns(fakeUser, nil)

	server := httptest.NewServer(http.HandlerFunc(facebookLoginHandler(formatter, userRepository, fbClient)))
	defer server.Close()

	validJSON := []byte("{\"id\":\"987\",\"access_token\":\"abc123\",\"signed_request\":\"abc123\",\"expires_in\":123}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(validJSON))
	if err != nil {
		t.Errorf("Error in creating POST request with facebookLoginHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to facebookLoginHandler: %v", err)
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected http.StatusOK, but got %v; response body: %v", res.StatusCode, string(payload))
	}

	if userRepository.FindByFacebookID(fakeUser.FacebookID) == nil {
		t.Error("Expected to find a user in the user repository with FacebookID of " + fakeUser.FacebookID)
	}

	var responseObject map[string]interface{}

	err = json.Unmarshal(payload, &responseObject)
	if err != nil {
		t.Errorf("Could not unmarshal payload into map[string]interface{} object")
	}

	if responseObject["Token"] == nil {
		t.Error("Response object should contain Token")
	}
}
