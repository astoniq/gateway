package router

import "net/http"

type Constructor func(handler http.Handler) http.Handler

type Router interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor)
	Any(path string, handler http.HandlerFunc, handlers ...Constructor)
	Get(path string, handler http.HandlerFunc, handlers ...Constructor)
	Post(path string, handler http.HandlerFunc, handlers ...Constructor)
	Put(path string, handler http.HandlerFunc, handlers ...Constructor)
	Delete(path string, handler http.HandlerFunc, handlers ...Constructor)
	Patch(path string, handler http.HandlerFunc, handlers ...Constructor)
	Head(path string, handler http.HandlerFunc, handlers ...Constructor)
	Options(path string, handler http.HandlerFunc, handlers ...Constructor)
	Trace(path string, handler http.HandlerFunc, handlers ...Constructor)
	Connect(path string, handler http.HandlerFunc, handlers ...Constructor)
	Group(path string)
	Use(handlers ...Constructor) Router
	Count() int
}

type Options struct {
	NotFoundHandler           http.HandlerFunc
	SafeAddRoutesWhileRunning bool
}

var DefaultOptions = Options{
	NotFoundHandler:           http.NotFound,
	SafeAddRoutesWhileRunning: true,
}
