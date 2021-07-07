package app

import "net/http"

// Response struct
type Response struct {
	Header      map[string]string
	StatusCode  int
	ContentType string
}

// Request struct
type Request struct {
	URI         string
	Method      string
	Header      map[string]string
	Response    Response
	Writer      http.ResponseWriter
	BaseRequest *http.Request

	// internal items...
	_getsData  map[string]interface{}
	_postsData map[string]interface{}
	_filesData map[string]interface{}
}

func (r *Request) Get(name string) interface{} {
	if val, ok := r._getsData[name]; ok {
		return val
	}
	return nil
}

func (r *Request) Post(name string) interface{} {
	if val, ok := r._postsData[name]; ok {
		return val
	}
	return nil
}

func (r *Request) File(name string) interface{} {
	if val, ok := r._filesData[name]; ok {
		return val
	}
	return nil
}
