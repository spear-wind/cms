package user

import (
	"errors"
	"fmt"

	"github.com/spear-wind/cms/security"
	"github.com/spear-wind/cms/validator"
	"gopkg.in/hlandau/passlib.v1"
)

type UserRepository interface {
	Add(user *User) (err error)
	Update(user *User) (err error)
	listUsers() (users []*User)
	getUser(id int64) (user *User, err error)
	Exists(user *User) bool
	FindByEmail(emailAddress string) (user *User)
	FindByVerificationCode(verificationCode string) (user *User)
	FindByFacebookID(facebookID string) (user *User)
}

type User struct {
	ID               int64  `json:"id"`
	FacebookID       string `json:"fb_id"`
	Email            string `json:"email"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Password         string `json:"password,omitempty"`
	hash             string
	Verified         bool   `json:"verified"`
	VerificationCode string `json:"-"`
}

type userListResponse struct {
	Total int    `json:"total"`
	Users []User `json:"users"`
}

func NewUser(ID int64, FirstName string, LastName string, Email string) *User {
	return &User{
		ID:        ID,
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
	}
}

func (user *User) String() string {
	return fmt.Sprintf("User{Id:%v, Verified:%v, Email:%v, Name:%v %v}", user.ID, user.Verified, user.Email, user.FirstName, user.LastName)
}

func (user *User) validate() (result validator.ValidationResult) {
	result = validator.NewValidationResult()

	if len(user.FirstName) == 0 {
		result.AddError("first_name", "First Name is required")
	}

	if len(user.LastName) == 0 {
		result.AddError("last_name", "Last Name is required")
	}

	if len(user.Email) == 0 {
		result.AddError("email", "Email is required")
	}

	if len(user.Password) == 0 {
		result.AddError("password", "Password is required")
	}

	return result
}

func (user *User) Register() (validator.ValidationResult, error) {
	result := user.validate()

	if result.HasErrors() {
		return result, nil
	}

	if user.Verified == true {
		return result, errors.New("This user has already been verified")
	}

	hash, err := passlib.Hash(user.Password)
	if err != nil {
		return result, fmt.Errorf("Failed to hash password: %v", err)
	}

	verificationCode, err := security.GenerateRandomString(24)
	if err != nil {
		return result, fmt.Errorf("Failed to generate verification code: %v", err)
	}

	user.hash = hash
	user.Password = ""
	user.VerificationCode = verificationCode

	return result, nil
}

func (user *User) Verify(verificationCode string) error {
	if user.Verified == true {
		return errors.New("This user has already been verified")
	}

	if user.VerificationCode != verificationCode {
		return errors.New("Invalid Verification Code")
	}

	user.VerificationCode = ""
	user.Verified = true

	return nil
}

func (user *User) Authenticate(password string) (bool, bool) {
	newHash, err := passlib.Verify(password, user.hash)
	if err != nil {
		fmt.Printf("passlib.Verify resulted in an error: %v", err)
		// incorrect password, malformed hash, etc.
		// either way, reject
		return false, false
	}

	hashUpdated := false
	// The context has decided, as per its policy, that
	// the hash which was used to validate the password
	// should be changed. It has upgraded the hash using
	// the verified password.
	if newHash != "" {
		user.Password = newHash
		hashUpdated = true
	}

	return true, hashUpdated
}
