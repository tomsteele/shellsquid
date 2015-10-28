package models

import (
	"net/http"
	"time"

	"github.com/mholt/binding"
	"github.com/nlf/boltons"
	"golang.org/x/crypto/bcrypt"
)

// User is a single user of the application.
type User struct {
	ID        string `json:"id"`
	Hash      string `json:"hash"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// NewUser creates and returns a new user object provided a username and password.
func NewUser(email string, password []byte) (*User, error) {
	now := time.Now().Unix()
	user := &User{
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	hash, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		return user, err
	}
	user.Hash = string(hash)
	return user, nil
}

// FindUserByEmail is a convenience function to locate a users record by a given email address.
func FindUserByEmail(db *boltons.DB, email string) (*User, error) {
	user := User{}
	users := []User{}
	if err := db.All(&users); err != nil {
		return &user, err
	}
	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}
	return &user, nil
}

// UserRequest is used for JSON binding when creating a new user.
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

func (u *UserRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if u.Email == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"email"},
			Message:    "email is required",
		})
	}
	if u.Password == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"password"},
			Message:    "password is required",
		})
	}
	return errs
}

// UserTokenRequest is used when doing a login to generate a new JWT token for the given user.
type UserTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserTokenRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

func (u *UserTokenRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if u.Email == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"email"},
			Message:    "email is required",
		})
	}
	if u.Password == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"password"},
			Message:    "password is required",
		})
	}
	return errs
}

// UserUpdateRequest is used to perform JSON binding when updating a user.
type UserUpdateRequest struct {
	Password string `json:"password"`
}

func (u *UserUpdateRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

func (u *UserUpdateRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	if u.Password == "" {
		errs = append(errs, binding.Error{
			FieldNames: []string{"password"},
			Message:    "password is required",
		})
	}
	return errs
}
