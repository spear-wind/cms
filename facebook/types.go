package facebook

import (
	"fmt"

	fb "github.com/huandu/facebook"
	"github.com/spear-wind/cms/user"
	"github.com/spear-wind/cms/validator"
)

type loginCommand struct {
	AccessToken string `json:"access_token"`
	//UNIX time when the token expires and needs to be renewed
	ExpiresIn     int64  `json:"expires_in"`
	SignedRequest string `json:"signed_request"`
	UserID        string `json:"id"`
}

func (cmd *loginCommand) validate() (result validator.ValidationResult) {
	result = validator.NewValidationResult()

	if len(cmd.AccessToken) == 0 {
		result.AddError("access_token", "Access Token is required")
	}

	if cmd.ExpiresIn == int64(0) {
		result.AddError("expires_in", "Expires In is required")
	}

	if len(cmd.SignedRequest) == 0 {
		result.AddError("signed_request", "Signed Request is required")
	}

	if len(cmd.UserID) == 0 {
		result.AddError("id", "id is required")
	}

	return result
}

func (cmd loginCommand) String() string {
	return fmt.Sprintf(
		"user id: %s\naccess token: %s\nexpires in: %d\nsigned request: %s\n",
		cmd.UserID,
		cmd.AccessToken,
		cmd.ExpiresIn,
		cmd.SignedRequest,
	)
}

type Client interface {
	getUser(cmd loginCommand) (*user.User, error)
}

type facebookClient struct {
	clientID     string
	clientSecret string
}

func NewClient(clientID string, clientSecret string) Client {
	return &facebookClient{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *facebookClient) getUser(cmd loginCommand) (*user.User, error) {
	user := &user.User{
		FacebookID: cmd.UserID,
	}

	//call /me to ensure that we can grab first name, last name, and email address
	res, err := fb.Get("/me", fb.Params{
		"fields":       "first_name,last_name,email",
		"access_token": cmd.AccessToken,
	})

	if err != nil {
		return nil, err
	}

	//assumes that we were able to retrieve info from Facebook API calls using the given access token
	user.Verified = true

	if firstName, ok := res["first_name"].(string); ok {
		user.FirstName = firstName
	}

	if lastName, ok := res["last_name"].(string); ok {
		user.LastName = lastName
	}

	if email, ok := res["email"].(string); ok {
		user.Email = email
	}

	return user, nil
}
