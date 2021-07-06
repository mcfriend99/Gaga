package main

import (
	"github.com/mcfriend99/gaga/app"
)

func main() {
	config := app.LoadConfig("config")
	g := app.Gaga{
		RouteGenerator: Routes,
		Config:         config,
	}
	g.Serve()
}
