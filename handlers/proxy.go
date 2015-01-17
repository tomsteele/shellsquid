package handlers

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/models"
)

func Proxy(server *app.App, isHttps bool) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		record, err := models.FindRecordByFQDN(server.DB, req.Host)
		if err != nil || record.ID == "" {
			server.Render.Data(w, http.StatusNotFound, nil)
			return
		}
		if record.Blacklist {
			server.Render.Data(w, http.StatusNotFound, nil)
			return
		}
		u, err := url.Parse(record.HandlerProtocol + "://" + record.HandlerHost + ":" + strconv.Itoa(record.HandlerPort))
		if err != nil {
			server.Render.Data(w, http.StatusNotFound, nil)
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
