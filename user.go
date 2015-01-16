package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/nlf/boltons"
)

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

type NewUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *NewUserRequest) FieldMap() binding.FieldMap {
	return binding.FieldMap{}
}
func (u *NewUserRequest) Validate(req *http.Request, errs binding.Errors) binding.Errors {
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

func UserCreateHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userReq := &NewUserRequest{}
		if err := binding.Bind(req, userReq); err.Handle(w) {
			return
		}
		existing, err := FindUserByEmail(app.DB, userReq.Email)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the user to the database"})
			log.Println(err)
			return
		}
		if existing.ID != "" {
			app.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "user email must be unique across the application"})
			return
		}
		user, err := NewUser(userReq.Email, []byte(userReq.Password))
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the user to the database"})
			log.Println(err)
			return
		}
		if err := app.DB.Save(user); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the user to the database"})
			log.Println(err)
			return
		}
		user.Hash = ""
		app.Render.JSON(w, http.StatusCreated, user)
	}
}

type UserTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserTokenRequest) FieldMap() binding.FieldMap {
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

func UserTokenHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userTokenReq := &UserTokenRequest{}
		if err := binding.Bind(req, userTokenReq); err.Handle(w) {
			return
		}
		user, err := FindUserByEmail(app.DB, userTokenReq.Email)
		if err != nil {
			log.Println(err)
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "error getting user from database"})
			return
		}
		if user.ID == "" {
			app.Render.JSON(w, http.StatusNotFound, map[string]string{"error": "invalid username or password"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(userTokenReq.Password)); err != nil {
			app.Render.JSON(w, http.StatusNotFound, map[string]string{"error": "invalid username or password"})
			return
		}
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["id"] = user.ID
		tokenString, err := token.SignedString(app.JWTSecret)
		if err != nil {
			log.Println(err)
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "error signing token"})
			return
		}
		app.Render.JSON(w, http.StatusCreated, map[string]string{"token": tokenString})
	}
}

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

func UserIndexHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		users := []User{}
		if err := app.DB.All(&users); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting users from the database"})
			log.Println(err)
			return
		}
		for _, u := range users {
			u.Hash = ""
		}
		app.Render.JSON(w, http.StatusOK, users)
	}
}

func UserShowHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		user := &User{ID: id}
		if ok, err := app.DB.Exists(user); err != nil || !ok {
			app.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if err := app.DB.Get(user); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting user from the database"})
			log.Println(err)
			return
		}
		user.Hash = ""
		app.Render.JSON(w, http.StatusOK, user)
	}
}

func UserDeleteHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		user := &User{ID: id}
		if ok, err := app.DB.Exists(user); err != nil || !ok {
			app.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if foundRecords, err := FindRecordsForOwner(app.DB, id); err != nil || len(foundRecords) > 0 {
			app.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "this user has records owned by them, remove or reassign before deleting"})
			return
		}
		if err := app.DB.Delete(user); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error deleting  the user from the database"})
			log.Println(err)
			return
		}
		app.Render.JSON(w, http.StatusNoContent, nil)
	}
}

type UserUpdateRequest struct {
	Password string `json:"password"`
}

func (u *UserUpdateRequest) FieldMap() binding.FieldMap {
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

func UserUpdateHandler(app *App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		user := &User{ID: id}
		if ok, err := app.DB.Exists(user); err != nil || !ok {
			app.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		updateReq := &UserUpdateRequest{}
		if err := binding.Bind(req, updateReq); err.Handle(w) {
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), 12)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the user"})
			log.Println(err)
			return
		}
		if err := app.DB.Update(user, map[string]interface{}{"Hash": string(hash)}); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the user"})
			log.Println(err)
			return
		}
		user.Hash = ""
		app.Render.JSON(w, http.StatusOK, user)
	}
}
