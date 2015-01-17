package handlers

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mholt/binding"
	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/models"
	"golang.org/x/crypto/bcrypt"
)

func UserToken(server *app.App) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userTokenReq := &models.UserTokenRequest{}
		if err := binding.Bind(req, userTokenReq); err.Handle(w) {
			return
		}
		user, err := models.FindUserByEmail(server.DB, userTokenReq.Email)
		if err != nil {
			log.Println(err)
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "error getting user from database"})
			return
		}
		if user.ID == "" {
			server.Render.JSON(w, http.StatusNotFound, map[string]string{"error": "invalid username or password"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(userTokenReq.Password)); err != nil {
			server.Render.JSON(w, http.StatusNotFound, map[string]string{"error": "invalid username or password"})
			return
		}
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["id"] = user.ID
		tokenString, err := token.SignedString(server.JWTSecret)
		if err != nil {
			log.Println(err)
			server.Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "error signing token"})
			return
		}
		server.Render.JSON(w, http.StatusCreated, map[string]string{"token": tokenString})
	}
}
