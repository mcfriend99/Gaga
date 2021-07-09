package app

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/mcfriend99/gaga/logger"
	"net/http"
	"regexp"
	"strconv"
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
		Routes:          make(map[string][]Route),
		_shouldCompress: g.Config.SEO.Compress,
	}

	// get user routes...
	g.RouteGenerator(&routing)

	// create request object
	contentType := r.Header.Get("Content-Type")

	request := Request{
		URI:    r.RequestURI,
		Method: r.Method,
		Header: make(map[string]string),
		Params: make(map[string]string),
		Response: Response{
			StatusCode:  http.StatusNotFound,
			ContentType: contentType,
		},
		Writer:      w,
		BaseRequest: r,
		_filesData:  make(map[string]interface{}),
		_postsData:  make(map[string]interface{}),
		_getsData:   make(map[string]interface{}),
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

	var routesToSearch []Route

	if routes, err := routing.Routes[r.Method]; err {
		routesToSearch = routes
	}

	if routes, err := routing.Routes[""]; err {
		routesToSearch = append(routesToSearch, routes...)
	}

	if len(routesToSearch) > 0 {
		for _, route := range routesToSearch {
			routeFound = route.Path == r.RequestURI

			// static routes should use a prefix with check.
			if route._isStatic {
				routeFound = strings.HasPrefix(r.RequestURI, route.Path)
			}

			if !routeFound {
				// do a regex check.
				// also, this is the only place where we can have route parameters.
				m := regexp.MustCompile(`/{([^}?]+)([?])?}/?`)
				path := m.ReplaceAllString(route.Path, "/?(?P<$1>[^/?#]+)$2/?")

				exp := regexp.MustCompile(fmt.Sprintf("^%s$", path))
				routeFound = exp.MatchString(r.RequestURI)

				if routeFound {
					match := exp.FindStringSubmatch(r.RequestURI)
					for i, name := range exp.SubexpNames() {
						if i != 0 && name != "" {
							// check param validator
							if v, e := route._paramValidators[name]; e && match[i] != "" {
								if k, e := regexp.MatchString(v, match[i]); e != nil || !k {
									routeFound = false
									break
								}
							}

							// use default param if the param is empty
							if v, e := route._paramDefaults[name]; e && match[i] == "" {
								request.Params[name] = v
							} else {
								request.Params[name] = match[i]
							}
						}
					}
				}
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

	contentLength := len(result)

	if _route != nil && _route._isStatic {
		contentLength, _ = strconv.Atoi(w.Header().Get("Content-Length"))
	} else {
		w.WriteHeader(request.Response.StatusCode)
	}

	if result != "" {
		writer := w.Write

		// only compress objects exceeding 128 byte
		if g.Config.SEO.Compress && len(result) > g.Config.SEO.CompressionThreshold {
			w.Header().Add("Vary", "Accept-Encoding")

			// prioritize gzip over deflate
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				gz := gzip.NewWriter(w)
				defer gz.Close()

				writer = gz.Write
			} else if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
				w.Header().Set("Content-Encoding", "deflate")

				fw, _ := flate.NewWriter(w, flate.BestCompression)
				defer fw.Close()

				writer = fw.Write
			}
		}

		if _, err := writer([]byte(result)); err != nil {
			logger.Error("Failed to write response:", err)
		}
	}

	// write response data...
	responseType := "text/plain; charset=utf-8"

	for key, value := range request.Response.Header {
		w.Header().Set(key, value)
		if key == "Content-Type" {
			responseType = value
		}
	}

	/*if _, e := request.Response.Header["Content-Type"]; !e {
		// @TODO: do automatic content-type detection here...
	}*/

	w.Header().Set("Content-Type", responseType)

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
		contentLength,
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
