package patmatch

import (
	"context"
	"net/http"
	"strings"
)

type Router struct {
	routes []*route
}

type route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{routes: make([]*route, 0)}
}

func (ro *Router) NewRoute(method, path string, handler http.HandlerFunc) {
	ro.routes = append(ro.routes, &route{Method: method, Pattern: path, Handler: handler})
}

func (ro *Router) Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range ro.routes {
		if match(r.URL.Path, route.Pattern) {

			// Debug code:
			// fmt.Println("Match Route", route.Method, route.Regex.String(), matches[0])

			if r.Method != route.Method {
				allow = append(allow, route.Method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, "Context")
			route.Handler(w, r.WithContext(ctx))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ","))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}

func match(path, pattern string) bool {
	for path != "" && pattern != "" {
		switch pattern[0] {
		case ':':
			slashPath := strings.IndexByte(path, '/')
			if slashPath < 0 {
				slashPath = len(path)
			}
			slashPattern := strings.IndexByte(pattern, '/')
			if slashPattern < 0 {
				slashPattern = len(pattern)
			}
			path = path[slashPath:]
			pattern = pattern[slashPattern:]
		case path[0]:
			path = path[1:]
			pattern = pattern[1:]
		default:
			return false
		}
	}
	return path == "" && pattern == ""
}
