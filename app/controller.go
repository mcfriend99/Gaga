package app

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
