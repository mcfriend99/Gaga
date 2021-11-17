package app

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mcfriend99/gaga/logger"
)

type Controller func(r *Request) string

// Resource interface is a form of controller that allows
// you to automatically create and bind CRUD operations to
// related methods in structs as well as automatically create
// related and useful routing.
type Resource interface {
	Create(r *Request)
	Update(r *Request)
	Delete(r *Request)
	View(r *Request)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s != nil && s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// Gzip Compression
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func StaticFileController(r *Request, prefix string, dir string, compress bool) string {
	// if strings.HasSuffix(prefix, "/") {
	// 	prefix = prefix[:len(prefix)-1]
	// }
	prefix = strings.TrimSuffix(prefix, "/")

	handler := http.StripPrefix(prefix, http.FileServer(neuteredFileSystem{http.Dir(dir)}))

	writer := r.Writer

	if compress {
		r.Writer.Header().Add("Vary", "Accept-Encoding")
		if strings.Contains(r.BaseRequest.Header.Get("Accept-Encoding"), "gzip") {
			r.Writer.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(r.Writer)
			defer gz.Close()

			writer = gzipResponseWriter{Writer: gz, ResponseWriter: r.Writer}
		} else if strings.Contains(r.BaseRequest.Header.Get("Accept-Encoding"), "deflate") {
			r.Writer.Header().Set("Content-Encoding", "deflate")
			fw, _ := flate.NewWriter(r.Writer, flate.BestCompression)
			defer fw.Close()

			writer = gzipResponseWriter{Writer: fw, ResponseWriter: r.Writer}
		}
	}

	if m, _ := regexp.MatchString("[.][a-zA-Z0-9]+$", r.URI); m {
		index := strings.LastIndex(r.URI, ".")
		ext := r.URI[index:len(r.URI)]
		logger.Infof("Checking mime type for %s based on extension %s...", r.URI, ext)

		responseType := mime.TypeByExtension(ext)
		logger.Infof("Static file mime type = %s", responseType)
		if responseType != "" {
			r.Response.Header["Content-Type"] = responseType
		}
	}

	handler.ServeHTTP(writer, r.BaseRequest)
	return ""
}
