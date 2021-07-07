package main

import (
	"github.com/mcfriend99/gaga/app"
	"github.com/mcfriend99/gaga/controller"
)

// Router is where you define routes for your application.
func Router(r *app.Routing) {
	// home page
	r.Get("/", controller.Home)

	// serving static page
	r.Static("/static", "./static")
}
