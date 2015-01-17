package app

import (
	"github.com/nlf/boltons"
	"github.com/unrolled/render"
)

// App is used by the server to pass around global data structures need by handlers.
type App struct {
	DB        *boltons.DB
	JWTSecret []byte
	Render    *render.Render
}
