package main

import (
	"github.com/mcfriend99/gaga/app"
)

func main() {
	config := app.LoadConfig()
	g := app.Gaga{
		RouteGenerator: Router,
		Config:         config,
	}
	g.Serve()
}
