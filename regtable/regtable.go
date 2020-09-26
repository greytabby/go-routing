package regtable

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Router struct {
	routes []*route
}

type route struct {
	Method  string
	Regex   *regexp.Regexp
	Handler http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{routes: make([]*route, 0)}
}

func (ro *Router) NewRoute(method, path string, handler http.HandlerFunc) {
	ro.routes = append(ro.routes, &route{Method: method, Regex: pathToRegex(path), Handler: handler})
}

// pathToRegex convert path to regex pattern
// example:
// /a/b/c -> /a/b/c
// /a/b/c/ -> /a/b/c
// /a/:b/c -> /a/([^/]+)/c
// /a/:b/:c -> /a/([^/]+)/([^/]+)
func pathToRegex(path string) *regexp.Regexp {
	path = strings.TrimRight(path, "/")
	subPaths := strings.Split(path, "/")
	for i, p := range subPaths {
		if strings.HasPrefix(p, ":") {
			subPaths[i] = "([^/]+)"
			continue
		}
	}
	pattern := strings.Join(subPaths, "/")
	if len(pattern) == 0 {
		pattern = "/"
	}
	pattern = "^" + pattern + "$"
	fmt.Println("Pattern:", pattern)
	return regexp.MustCompile(pattern)
}

func (ro *Router) Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range ro.routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {

			// Debug code:
			// fmt.Println("Match Route", route.Method, route.Regex.String(), matches[0])

			if r.Method != route.Method {
				allow = append(allow, route.Method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
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
