package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/models"
)

// SetUserContext takes a JWT token from the mux context and looks up the user by id. The user
// is then set into the same context.
func SetUserContext(server *app.App) func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		token := context.Get(req, "user").(*jwt.Token)
		id := token.Claims["id"].(string)
		user := &models.User{ID: id}
		if ok, err := server.DB.Exists(user); err != nil || !ok {
			server.Render.JSON(w, http.StatusUnauthorized, nil)
			return
		}
		if err := server.DB.Get(user); err != nil {
			server.Render.JSON(w, http.StatusUnauthorized, nil)
			return
		}
		context.Set(req, "user", user)
		next(w, req)
	}
}
