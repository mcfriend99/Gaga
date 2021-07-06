package main

import "github.com/mcfriend99/gaga/app"

// Routes is where you define routes for your application.
func Routes(r *app.Routing) {
	// home page
	r.Get("/", Home)
}
