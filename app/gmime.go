package app

import "mime"

func InitGagaMimes() {
	// setup mime types for static files...
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".html", "text/html")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".jpg", "images/jpeg")
	mime.AddExtensionType(".png", "images/png")
	mime.AddExtensionType(".gif", "images/gif")
}
