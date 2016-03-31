package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/dave-malone/email"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func TestCreateUserHandlerResponseToInvalidData(t *testing.T) {
	client := &http.Client{}
	email.NewSender = email.NewNoopSender
	repo := newInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createUserHandler(formatter, repo)))
	defer server.Close()

	invalidBody := []byte("not even json")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(invalidBody))
	if err != nil {
		t.Errorf("Error in creating POST request with CreateUserHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to CreateUserHandler: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending invalid Body should result in a bad request from the server")
	}
}

func TestCreateUserHandlerResponseToBadJson(t *testing.T) {
	client := &http.Client{}
	email.NewSender = email.NewNoopSender
	repo := newInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createUserHandler(formatter, repo)))
	defer server.Close()

	badJSON := []byte("{\"test\":\"bad json! bad!\"}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(badJSON))
	if err != nil {
		t.Errorf("Error in creating POST request with createUserHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createUserHandler: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending bad JSON should result in a bad request from the server")
	}
}

func TestCreateUserHandler(t *testing.T) {
	client := &http.Client{}
	email.NewSender = email.NewNoopSender
	repo := newInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createUserHandler(formatter, repo)))
	defer server.Close()

	body := []byte("{\"first_name\":\"john\", \"last_name\":\"doe\", \"email\":\"john@doe.com\", \"password\":\"p@$$w3Rd\"}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createUserHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createUserHandler: %v", err)
	}

	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, received %s", res.Status)
	}

	loc, headerOk := res.Header["Location"]

	if !headerOk {
		t.Error("Location header is not set")
	} else {
		if !strings.Contains(loc[0], "/user/") {
			t.Errorf("Location header should contain '/user/'")
		}
		loc, headerOk := res.Header["Location"]

		if !headerOk {
			t.Error("Location header is not set")
		} else {
			if !strings.Contains(loc[0], "/user/") {
				t.Errorf("Location header should contain '/user/'")
			}
			if loc[0] != "/user/0" {
				t.Errorf("Expected '/user/0' but was %v", loc[0])
			}
		}
	}
	var user User
	err = json.Unmarshal(payload, &user)
	if err != nil {
		t.Errorf("Could not unmarshal payload into User object")
	}

	if user.ID != int64(0) || !strings.Contains(loc[0], strconv.FormatInt(user.ID, 10)) {
		t.Error("user.ID does not match Location header")
	}

	users := repo.listUsers()
	if len(users) != 1 {
		t.Errorf("Expected user repo to have exactly 1 user, but there were %d", len(users))
	}

	if users[0].FirstName != "john" {
		t.Errorf("Repo user first name should be 'john', but was %s", users[0].FirstName)
	}

	if users[0].LastName != "doe" {
		t.Errorf("Repo user last name should be 'doe', but was %s", users[0].FirstName)
	}

	if users[0].Email != "john@doe.com" {
		t.Errorf("Repo user email should be 'john@doe.com', but was %s", users[0].FirstName)
	}

}

func TestGetUserListReturnsEmptyArrayForNoUsers(t *testing.T) {
	client := &http.Client{}
	email.NewSender = email.NewNoopSender
	repo := newInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(getUserListHandler(formatter, repo)))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Error("Errored when sending request to the server", err)
		return
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response from server", err)
	}

	var userListResponse userListResponse

	err = json.Unmarshal(payload, &userListResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into data struct: %v", err)
	}

	if len(userListResponse.Users) != 0 {
		t.Errorf("Expected an empty list of user responses, got %d", len(userListResponse.Users))
	}
}

func TestGetUserListReturnsWhatsInRepository(t *testing.T) {
	client := &http.Client{}
	email.NewSender = email.NewNoopSender
	repo := newInMemoryRepository()
	repo.add(*newUser(-1, "John", "Doe", "john@doe.com"))
	repo.add(*newUser(-1, "Jane", "Doe", "jane@doe.com"))
	repo.add(*newUser(-1, "Baby", "Doe", "baby@doe.com"))
	server := httptest.NewServer(http.HandlerFunc(getUserListHandler(formatter, repo)))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)
	if err != nil {
		t.Error("Errored when sending request to the server", err)
		return
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response from server", err)
	}

	var userListResponse userListResponse

	err = json.Unmarshal(payload, &userListResponse)
	if err != nil {
		t.Errorf("Could not unmarshal payload into data struct: %v", err)
	}

	if len(userListResponse.Users) != 3 {
		t.Errorf("Expected exactly three users in the user response, but got %d", len(userListResponse.Users))
	}
}

func MakeTestServer() *negroni.Negroni {
	server := negroni.New() // don't need all the middleware here or logging.
	router := mux.NewRouter()
	InitRoutes(router, formatter)
	server.UseHandler(router)
	return server
}
