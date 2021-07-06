package main

import "github.com/mcfriend99/gaga/gaga"

// Routes is where you define routes for your application.
func Routes(r *gaga.Routing) {
	// home page
	r.Get("/", Home)
}
