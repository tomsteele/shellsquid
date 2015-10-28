package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmcvetta/randutil"
	"github.com/miekg/dns"
	"github.com/nlf/boltons"
	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/config"
	"github.com/tomsteele/shellsquid/handlers"
	"github.com/tomsteele/shellsquid/middleware"
	"github.com/tomsteele/shellsquid/models"
	"github.com/unrolled/render"
)

const version = "2.2.0"

func main() {
	conf, err := config.New("./config.json")
	if err != nil {
		log.Fatalf("Error parsing confuration file: %s", err.Error())
	}
	if conf.JWTKey == "" {
		log.Fatalf("jwt_secret in config.json is not set, please set this to a random value")
	}

	db, err := boltons.Open(conf.BoltDBFile, 0600, nil)
	if err != nil {
		log.Fatalf("Error opening db: %s", err.Error())
	}
	defer db.Close()

	keys, err := db.Keys(models.User{})
	if err != nil {
		log.Fatalf("Error getting keys from db: %s", err.Error())
	}
	if len(keys) == 0 {
		log.Println("No users found creating admin@localhost user with random password")
		random, err := randutil.AlphaString(10)
		if err != nil {
			log.Fatalf("Error generating password: %s", err.Error())
		}
		firstUser, err := models.NewUser("admin@localhost", []byte(random))
		if err != nil {
		}
		if err := db.Save(firstUser); err != nil {
			log.Fatalf("Error saving first user to db: %s", err.Error())
		}
		log.Printf("admin@localhost password set to %s", random)
	}

	serverApp := &app.App{
		DB:        db,
		JWTSecret: []byte(conf.JWTKey),
		Render:    render.New(),
		Config:    conf,
	}

	if conf.Proxy.SSL.Enabled {
		sslMux := http.NewServeMux()
		sslMux.HandleFunc("/", handlers.Proxy(serverApp, true))
		sslRecovery := negroni.NewRecovery()
		sslRecovery.PrintStack = false
		sslProxy := negroni.New(sslRecovery)
		sslProxy.UseHandler(sslMux)
		go func() {
			log.Fatal(http.ListenAndServeTLS(conf.Proxy.SSL.Listener, conf.Proxy.SSL.Cert, conf.Proxy.SSL.Key, sslProxy))
		}()
	}

	if conf.Proxy.HTTP.Enabled {
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/", handlers.Proxy(serverApp, false))
		httpRecovery := negroni.NewRecovery()
		httpRecovery.PrintStack = false
		httpProxy := negroni.New(httpRecovery)
		httpProxy.UseHandler(httpMux)
		go func() {
			log.Fatal(http.ListenAndServe(conf.Proxy.HTTP.Listener, httpProxy))
		}()
	}

	if conf.Proxy.DNS.Enabled {
		tcpDNServer := &dns.Server{Addr: conf.Proxy.DNS.Listener, Net: "udp"}
		udpDNServer := &dns.Server{Addr: conf.Proxy.DNS.Listener, Net: "tcp"}

		dns.HandleFunc(".", handlers.ProxyDNS(serverApp))

		go func() {
			log.Fatal(tcpDNServer.ListenAndServe())
		}()

		go func() {
			log.Fatal(udpDNServer.ListenAndServe())
		}()
	}

	r := mux.NewRouter()
	api := mux.NewRouter()
	r.HandleFunc("/api/token", handlers.UserToken(serverApp)).Methods("POST")

	api.HandleFunc("/api/users", handlers.CreateUser(serverApp)).Methods("POST")
	api.HandleFunc("/api/users", handlers.IndexUser(serverApp)).Methods("GET")
	api.HandleFunc("/api/users/{id}", handlers.ShowUser(serverApp)).Methods("GET")
	api.HandleFunc("/api/users/{id}", handlers.DeleteUser(serverApp)).Methods("DELETE")
	api.HandleFunc("/api/users/{id}", handlers.UpdateUser(serverApp)).Methods("PUT")
	api.HandleFunc("/api/records", handlers.CreateRecord(serverApp)).Methods("POST")
	api.HandleFunc("/api/records", handlers.IndexRecord(serverApp)).Methods("GET")
	api.HandleFunc("/api/records/{id}", handlers.ShowRecord(serverApp)).Methods("GET")
	api.HandleFunc("/api/records/{id}", handlers.DeleteRecord(serverApp)).Methods("DELETE")
	api.HandleFunc("/api/records/{id}", handlers.UpdateRecord(serverApp)).Methods("PUT")
	api.HandleFunc("/api/info", handlers.Info(serverApp, version)).Methods("GET")

	r.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(middleware.JWTAuth(serverApp)),
		negroni.HandlerFunc(middleware.SetUserContext(serverApp)),
		negroni.Wrap(api),
	))

	server := negroni.New(
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("static")),
		negroni.NewRecovery(),
	)

	server.UseHandler(r)
	log.Fatal(http.ListenAndServeTLS(conf.Admin.Listener, conf.Admin.Cert, conf.Admin.Key, server))

}
