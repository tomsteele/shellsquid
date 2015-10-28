package handlers

import (
	"net/http"

	"github.com/tomsteele/shellsquid/app"
)

// Infos struct holds version and proxy info.
type Infos struct {
	Version string `json:"version"`
	Proxy   struct {
		DNS struct {
			Enabled  bool   `json:"enabled"`
			Listener string `json:"listener"`
		} `json:"dns"`
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
func Info(server *app.App, version string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		info := &Infos{}
		info.Version = version
		info.Proxy.SSL.Enabled = server.Config.Proxy.SSL.Enabled
		info.Proxy.SSL.Listener = server.Config.Proxy.SSL.Listener
		info.Proxy.HTTP.Enabled = server.Config.Proxy.HTTP.Enabled
		info.Proxy.HTTP.Listener = server.Config.Proxy.HTTP.Listener
		info.Proxy.DNS.Enabled = server.Config.Proxy.DNS.Enabled
		info.Proxy.DNS.Listener = server.Config.Proxy.DNS.Listener
		server.Render.JSON(w, http.StatusOK, info)
	}
}
