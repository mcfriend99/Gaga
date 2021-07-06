package gaga

import (
	"fmt"
	"log"
	"net/http"
)

// Gaga main struct
type Gaga struct {
	Config          *Config
	RouteGenerator  func(*Routing)
	NotFoundHandler func(*Request) string
}

func (g Gaga) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// initialize routing
	routing := Routing{
		Routes: make(map[string][]Route, 0),
	}

	// get user routes...
	g.RouteGenerator(&routing)

	// create request object
	contentType := r.Header.Get("Content-Type")

	request := Request{
		URI:    r.RequestURI,
		Method: r.Method,
		Header: make(map[string]string, 0),
		Response: Response{
			StatusCode:  http.StatusNotFound,
			ContentType: contentType,
		},
		_filesData: make(map[string]interface{}, 0),
		_postsData: make(map[string]interface{}, 0),
		_getsData:  make(map[string]interface{}, 0),
	}

	// populate request bodies...
	for s := range r.Header {
		request.Header[s] = r.Header.Get(s)
	}
	for s := range r.Form {
		request._getsData[s] = r.Form.Get(s)
	}
	for s := range r.PostForm {
		request._postsData[s] = r.PostForm.Get(s)
	}

	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		for s := range r.MultipartForm.File {
			request._filesData[s] = r.MultipartForm.File[s]
		}
	}

	// match route...
	result := ""
	routeFound := false

	if routes, err := routing.Routes[r.Method]; err {
		for _, route := range routes {
			if route.Path == r.RequestURI {
				routeFound = true
				request.Response.StatusCode = http.StatusOK
				if route.Controller != nil {
					result = route.Controller(&request)
				}
				break
			}
		}
	}

	// handle not found...
	if !routeFound {
		// handle with 404 handler here...
		// if non is given, do the default below.
		if g.NotFoundHandler != nil {
			request.Response.StatusCode = http.StatusOK
			result = g.NotFoundHandler(&request)
		} else {
			result = "404 Not Found"
		}
	}

	// write response data...
	responseType := "text/html"
	if v, e := request.Response.Header["Content-Type"]; e {
		responseType = v
	} else {
		// @TODO: do automatic content-type detection here...
	}
	w.WriteHeader(request.Response.StatusCode)
	w.Write([]byte(result))

	// log request
	log.Printf(`%s "%s %s %s" %d %d "%s - %s" "%s"`,
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Proto,
		request.Response.StatusCode,
		len(result),
		contentType,
		responseType,
		r.UserAgent(),
	)
}

func (g *Gaga) Serve() {
	listen := fmt.Sprintf("%s:%d", g.Config.Server.ListenOn, g.Config.Server.Port)

	if !g.Config.Server.Secure {
		log.Printf("Started serving HTTP on http://%s\n", listen)
	} else {
		log.Printf("Started serving HTTPS on https://%s\n", listen)
	}

	if g.Config.Server.Secure {
		http.ListenAndServeTLS(listen,
			fmt.Sprintf("ssl/%s", g.Config.Server.TLSCertificateFile),
			fmt.Sprintf("ssl/%s", g.Config.Server.TLSKeyFile),
			g,
		)
	} else {
		http.ListenAndServe(listen, g)
	}
}
