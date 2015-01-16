package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func ProxyHandler(app *App, isHttps bool) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		record, err := FindRecordByFQDN(app.DB, req.Host)
		if err != nil || record.ID == "" {
			app.Render.Data(w, http.StatusNotFound, nil)
			return
		}
		if record.Blacklist {
			app.Render.Data(w, http.StatusNotFound, nil)
			return
		}
		u, err := url.Parse(record.HandlerProtocol + "://" + record.HandlerHost + ":" + strconv.Itoa(record.HandlerPort))
		if err != nil {
			app.Render.Data(w, http.StatusNotFound, nil)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(u)
		if isHttps {
			proxy.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
		proxy.ServeHTTP(w, req)
	}
}
