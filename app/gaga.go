package app

import (
	"fmt"
	"github.com/mcfriend99/gaga/logger"
	"net/http"
	"strings"
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
		Writer:      w,
		BaseRequest: r,
		_filesData:  make(map[string]interface{}, 0),
		_postsData:  make(map[string]interface{}, 0),
		_getsData:   make(map[string]interface{}, 0),
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
	var _route *Route = nil

	if routes, err := routing.Routes[r.Method]; err {
		for _, route := range routes {
			routeFound = route.Path == r.RequestURI

			// static routes should use a prefix with check.
			if route.IsStatic {
				routeFound = strings.HasPrefix(r.RequestURI, route.Path)
			}

			if routeFound {
				_route = &route
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
			// @TODO: use beautiful template based 404 page.
			result = "404 page not found"
		}
	}

	// write response data...
	responseType := "text/html"
	if v, e := request.Response.Header["Content-Type"]; e {
		responseType = v
	} else {
		// @TODO: do automatic content-type detection here...
	}

	if _route != nil && _route.IsStatic {
		// do nothing
	} else {
		w.WriteHeader(request.Response.StatusCode)
	}
	if result != "" {
		w.Write([]byte(result))
	}

	// logs request
	logMethod := logger.Infof
	if request.Response.StatusCode < 200 || request.Response.StatusCode >= 399 {
		logMethod = logger.Warnf
	}
	logMethod(`%s "%s %s %s" %d %d "%s - %s" "%s"`,
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

func (g *Gaga) setupLogging() {

	engine := logger.LogDestBoth

	if g.Config.Log.Engine == "file" {
		engine = logger.LogDestFile
	} else if g.Config.Log.Engine == "console" {
		engine = logger.LogDestConsole
	}

	var flag logger.ControlFlag
	level := logger.LogLevelInfo

	if g.Config.Log.ShowSource {
		flag = logger.ControlFlag(int(flag) | int(logger.ControlFlagLogLineNum) | int(logger.ControlFlagLogFuncName))
	}

	if g.Config.Log.Level == "info" {
		level = logger.LogLevelInfo
	} else if g.Config.Log.Level == "error" {
		level = logger.LogLevelError
	} else if g.Config.Log.Level == "warn" {
		level = logger.LogLevelWarn
	} else if g.Config.Log.Level == "fatal" {
		level = logger.LogLevelFatal
	} else if g.Config.Log.Level == "panic" {
		level = logger.LogLevelPanic
	} else if g.Config.Log.Level == "trace" {
		level = logger.LogLevelTrace
	}

	logger.Init(&logger.Config{
		LogDir:          g.Config.Log.Path,
		LogFileMaxSize:  200,
		LogFileMaxNum:   500,
		LogFileNumToDel: 50,
		LogLevel:        level,
		LogDest:         engine,
		Flag:            flag,
	})

	logger.Info("File logging initialized...")
}

func (g *Gaga) Serve() {
	g.setupLogging()

	listen := fmt.Sprintf("%s:%d", g.Config.Server.ListenOn, g.Config.Server.Port)

	if !g.Config.Server.Secure {
		logger.Infof("Started serving HTTP on http://%s\n", listen)
	} else {
		logger.Infof("Started serving HTTPS on https://%s\n", listen)
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
