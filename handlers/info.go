package handlers

import (
	"net/http"

	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/config"
)

// Infos struct holds version and proxy info.
type Infos struct {
	Version string `json:"version"`
	Proxy   struct {
		SSL struct {
			Enabled  bool   `json:"enabled"`
			Listener string `json:"listener"`
		} `json:"ssl"`
		HTTP struct {
			Enabled  bool   `json:"enabled"`
			Listener string `json:"listener"`
		} `json:"http"`
	} `json:"proxy"`
}

// Info returns some useful information for the client.
func Info(server *app.App, version string, conf *config.Config) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		info := &Infos{}
		info.Version = version
		info.Proxy.SSL.Enabled = conf.Proxy.SSL.Enabled
		info.Proxy.SSL.Listener = conf.Proxy.SSL.Listener
		info.Proxy.HTTP.Enabled = conf.Proxy.SSL.Enabled
		info.Proxy.HTTP.Listener = conf.Proxy.SSL.Listener
		server.Render.JSON(w, http.StatusOK, info)
	}
}
