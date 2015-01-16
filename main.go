package main

import (
	"log"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jmcvetta/randutil"
	"github.com/nlf/boltons"
	"github.com/unrolled/render"
)

func main() {
	config, err := ParseConfig("./config.json")
	if err != nil {
		log.Fatalf("Error parsing configuration file: %s", err.Error())
	}
	db, err := boltons.Open(config.BoltDBFile, 0600, nil)
	if err != nil {
		log.Fatal("Error opening db: %s", err.Error())
	}
	defer db.Close()

	keys, err := db.Keys(User{})
	if err != nil {
		log.Fatal("Error getting keys from db: %s", err.Error())
	}
	if len(keys) == 0 {
		log.Println("No users found creating admin@localhost user with random password")
		random, err := randutil.AlphaString(10)
		if err != nil {
			log.Fatalf("Error generating password: %s", err.Error())
		}
		firstUser, err := NewUser("admin@localhost", []byte(random))
		if err != nil {
		}
		if err := db.Save(firstUser); err != nil {
			log.Fatalf("Error saving first user to db: %s", err.Error())
		}
		log.Printf("admin@localhost password set to %s", random)
	}

	app := &App{
		DB:        db,
		JWTSecret: []byte(config.JWTKey),
		Render:    render.New(),
	}

	if config.Proxy.SSL.Enabled {
		sslMux := http.NewServeMux()
		sslMux.HandleFunc("/", ProxyHandler(app, false))
		sslProxy := negroni.Classic()
		sslProxy.UseHandler(sslMux)
		go func() {
			log.Fatal(http.ListenAndServeTLS(config.Proxy.SSL.Listener, config.Proxy.SSL.Cert, config.Proxy.SSL.Key, sslProxy))
		}()
	}

	if config.Proxy.HTTP.Enabled {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", ProxyHandler(app, false))
		httpProxy := negroni.Classic()
		httpProxy.UseHandler(httpMux)
		go func() {
			log.Fatal(http.ListenAndServe(config.Proxy.HTTP.Listener, httpProxy))
		}()
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return app.JWTSecret, nil
		},
	})

	r := mux.NewRouter()
	recordRouter := mux.NewRouter()
	userRouter := mux.NewRouter()

	r.HandleFunc("/api/token", UserTokenHandler(app)).Methods("POST")
	userRouter.HandleFunc("/api/users", UserCreateHandler(app)).Methods("POST")
	userRouter.HandleFunc("/api/users", UserIndexHandler(app)).Methods("GET")
	userRouter.HandleFunc("/api/users/{id}", UserShowHandler(app)).Methods("GET")
	userRouter.HandleFunc("/api/users/{id}", UserDeleteHandler(app)).Methods("DELETE")
	userRouter.HandleFunc("/api/users/{id}", UserUpdateHandler(app)).Methods("PUT")

	recordRouter.HandleFunc("/api/records", RecordCreateHandler(app)).Methods("POST")
	recordRouter.HandleFunc("/api/records", RecordIndexHandler(app)).Methods("GET")
	recordRouter.HandleFunc("/api/records/{id}", RecordShowHandler(app)).Methods("GET")
	recordRouter.HandleFunc("/api/records/{id}", RecordDeleteHandler(app)).Methods("DELETE")
	recordRouter.HandleFunc("/api/records/{id}", RecordUpdateHandler(app)).Methods("PUT")

	r.PathPrefix("/api/records").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(SetUserContext(app)),
		negroni.Wrap(recordRouter),
	))
	r.PathPrefix("/api/users").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(SetUserContext(app)),
		negroni.Wrap(userRouter),
	))

	server := negroni.Classic()
	server.UseHandler(r)
	log.Fatal(http.ListenAndServeTLS(config.Admin.Listener, config.Admin.Cert, config.Admin.Key, server))
}
