package router

import "net/http"

type Router interface {
	NewRoute(method, path string, handler http.Handler)
	Serve(w http.ResponseWriter, r *http.Request)
}
