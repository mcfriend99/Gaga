package app

import (
	"net/http"
	"path/filepath"
	"strings"
)

type Controller func(*Request) string

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

	s, err := f.Stat()
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

func StaticFileController(r *Request, prefix string, path string) string {
	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[:len(prefix)-1]
	}

	handler := http.StripPrefix(prefix, http.FileServer(neuteredFileSystem{http.Dir(path)}))
	handler.ServeHTTP(r.Writer, r.BaseRequest)
	return ""
}
