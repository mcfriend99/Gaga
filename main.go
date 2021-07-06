package main

import (
	"github.com/mcfriend99/gaga/gaga"
)

func main() {
	config := gaga.LoadConfig("config")
	g := gaga.Gaga{
		RouteGenerator: Routes,
		Config:         config,
	}
	g.Serve()
}
