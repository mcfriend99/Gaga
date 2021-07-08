package app

// Route struct
type Route struct {
	Path             string
	Controller       Controller
	_isStatic        bool
	_paramValidators map[string]string
}

// Where allows specifying a pattern that a named param in a
//  route must match to be valid.
//
//  Example:
//
//  r.Route("/{id}", controller.User).Where("id", `\d+`)
func (r *Route) Where(name string, test string) *Route {
	r._paramValidators[name] = test
	return r
}

// Routing struct
type Routing struct {
	Routes          map[string][]Route
	_shouldCompress bool
}

// CreateRoute allows you to create a route for any HTTP method
// bound to a given path.
func (r *Routing) CreateRoute(method string, path string, controller Controller) *Route {
	if r.Routes[method] == nil {
		r.Routes[method] = make([]Route, 0)
	}

	route := Route{
		Path:             path,
		Controller:       controller,
		_paramValidators: make(map[string]string),
	}

	r.Routes[method] = append(r.Routes[method], route)

	return &route
}

// Get routes an HTTP GET request with a request URI matching
// the given path to the given controller.
func (r *Routing) Get(path string, controller Controller) *Route {
	return r.CreateRoute("GET", path, controller)
}

// Post routes an HTTP POSt request with a request URI matching
// the given path to the given controller.
func (r *Routing) Post(path string, controller Controller) *Route {
	return r.CreateRoute("POST", path, controller)
}

// Put routes an HTTP PUT request with a request URI matching
// the given path to the given controller.
func (r *Routing) Put(path string, controller Controller) *Route {
	return r.CreateRoute("PUT", path, controller)
}

// Delete routes an HTTP DELETE request with a request URI matching
// the given path to the given controller.
func (r *Routing) Delete(path string, controller Controller) *Route {
	return r.CreateRoute("DELETE", path, controller)
}

// Patch routes an HTTP PATCH request with a request URI matching
// the given path to the given controller.
func (r *Routing) Patch(path string, controller Controller) *Route {
	return r.CreateRoute("PATCH", path, controller)
}

// Trace routes an HTTP TRACE request with a request URI matching
// the given path to the given controller.
func (r *Routing) Trace(path string, controller Controller) *Route {
	return r.CreateRoute("TRACE", path, controller)
}

// Options routes an HTTP OPTIONS request with a request URI matching
// the given path to the given controller.
func (r *Routing) Options(path string, controller Controller) *Route {
	return r.CreateRoute("OPTIONS", path, controller)
}

// Head routes an HTTP HEAD request with a request URI matching
// the given path to the given controller.
func (r *Routing) Head(path string, controller Controller) *Route {
	return r.CreateRoute("OPTIONS", path, controller)
}

// Static routes a request for static files matching the path to
// the specified directory
func (r *Routing) Static(path string, dir string) {
	if r.Routes["GET"] == nil {
		r.Routes["GET"] = make([]Route, 0)
	}

	r.Routes["GET"] = append(r.Routes["GET"], Route{
		Path: path,
		Controller: func(h *Request) string {
			return StaticFileController(h, path, dir, r._shouldCompress)
		},
		_isStatic:        true,
		_paramValidators: make(map[string]string),
	})
}

// Any allows creation of a route that's bound to all/any request method.
func (r *Routing) Any(path string, controller Controller) *Route {
	return r.CreateRoute("", path, controller)
}
