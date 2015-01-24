package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/models"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser is a http handler function to creation a new user.
func CreateUser(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userReq := &models.UserRequest{}
		if err := binding.Bind(req, userReq); err.Handle(w) {
			return
		}
		existing, err := models.FindUserByEmail(server.DB, userReq.Email)
		if err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the user to the database"})
			log.Println(err)
			return
		}
		if existing.ID != "" {
			server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "user email must be unique across the application"})
			return
		}
		user, err := models.NewUser(userReq.Email, []byte(userReq.Password))
		if err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the user to the database"})
			log.Println(err)
			return
		}
		if err := server.DB.Save(user); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error saving the user to the database"})
			log.Println(err)
			return
		}
		user.Hash = ""
		server.Render.JSON(w, http.StatusCreated, user)
	}
}

// IndexUser returns a list of all users.
func IndexUser(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		users := []models.User{}
		if err := server.DB.All(&users); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting users from the database"})
			log.Println(err)
			return
		}
		for _, u := range users {
			u.Hash = ""
		}
		server.Render.JSON(w, http.StatusOK, users)
	}
}

// ShowUser returns a single user provided by an id mux parameter.
func ShowUser(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		user := &models.User{ID: id}
		if ok, err := server.DB.Exists(user); err != nil || !ok {
			server.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if err := server.DB.Get(user); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error getting user from the database"})
			log.Println(err)
			return
		}
		user.Hash = ""
		server.Render.JSON(w, http.StatusOK, user)
	}
}

// DeleteUser deletes a single user provided by an id mux parameter.
func DeleteUser(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		user := &models.User{ID: id}
		if ok, err := server.DB.Exists(user); err != nil || !ok {
			server.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		if foundRecords, err := models.FindRecordsForOwner(server.DB, id); err != nil || len(foundRecords) > 0 {
			server.Render.JSON(w, http.StatusBadRequest, map[string]string{"error": "this user has records owned by them, remove or reassign before deleting"})
			return
		}
		if err := server.DB.Delete(user); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error deleting the user from the database"})
			log.Println(err)
			return
		}
		server.Render.JSON(w, http.StatusNoContent, nil)
	}
}

// UpdateUser updates selected fields of a user provided by an id mux parameter.
func UpdateUser(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		user := &models.User{ID: id}
		if ok, err := server.DB.Exists(user); err != nil || !ok {
			server.Render.JSON(w, http.StatusNotFound, nil)
			return
		}
		updateReq := &models.UserUpdateRequest{}
		if err := binding.Bind(req, updateReq); err.Handle(w) {
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), 12)
		if err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the user"})
			log.Println(err)
			return
		}
		if err := server.DB.Update(user, map[string]interface{}{"Hash": string(hash)}); err != nil {
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "there was an error updating the user"})
			log.Println(err)
			return
		}
		user.Hash = ""
		server.Render.JSON(w, http.StatusOK, user)
	}
}
