package handlers

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/tomsteele/shellsquid/app"
	"github.com/tomsteele/shellsquid/models"
)

func hostname(host string) string {
	parts := strings.SplitN(host, ":", 2)
	return parts[0]
}

func ProxyDNS(server *app.App) func(w dns.ResponseWriter, req *dns.Msg) {
	return func(w dns.ResponseWriter, req *dns.Msg) {
		if len(req.Question) == 0 {
			dns.HandleFailed(w, req)
			return
		}
		name := req.Question[0].Name
		record, err := models.FindRecordBySubOfFQDN(server.DB, name)
		if err != nil || record.ID == "" {
			dns.HandleFailed(w, req)
			return
		}
		if record.Blacklist {
			dns.HandleFailed(w, req)
			return
		}
		transport := "udp"
		if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
			transport = "tcp"
		}
		c := &dns.Client{Net: transport}
		resp, _, err := c.Exchange(req, record.HandlerHost+":"+strconv.Itoa(record.HandlerPort))
		if err != nil {
			dns.HandleFailed(w, req)
			return
		}
		w.WriteMsg(resp)
	}
}

func Proxy(server *app.App, isHttps bool) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		record, err := models.FindRecordByFQDN(server.DB, hostname(req.Host))
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
