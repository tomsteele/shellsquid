package main

import (
	"github.com/nlf/boltons"
	"github.com/unrolled/render"
)

type App struct {
	DB        *boltons.DB
	JWTSecret []byte
	Render    *render.Render
}
