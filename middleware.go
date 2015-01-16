package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func SetUserContext(app *App) func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		token := context.Get(req, "user").(*jwt.Token)
		id := token.Claims["id"].(string)
		user := &User{ID: id}
		if err := app.DB.Get(user); err != nil {
			app.Render.JSON(w, http.StatusUnauthorized, nil)
			return
		}
		context.Set(req, "user", user)
		next(w, req)
	}
}
