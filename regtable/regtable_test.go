package regtable_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/greytabby/go-routing/regtable"
	"github.com/stretchr/testify/assert"
)

func TestRouting(t *testing.T) {
	router := regtable.NewRouter()
	router.NewRoute(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/")
	}))
	router.NewRoute(http.MethodGet, "/api/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/api/test")
	}))
	router.NewRoute(http.MethodGet, "/api/test/:slug", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/api/test/:slug")
	}))
	router.NewRoute(http.MethodGet, "/test/:slug", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/test/:slug")
	}))
	router.NewRoute(http.MethodGet, "/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Get: /test")
	}))
	router.NewRoute(http.MethodPost, "/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Post: /test")
	}))

	cases := []struct {
		name   string
		method string
		path   string
		want   string
	}{
		{name: "/", method: http.MethodGet, path: "/", want: "/"},
		{name: "/api/test", method: http.MethodGet, path: "/api/test", want: "/api/test"},
		{name: "/api/test/:slug", method: http.MethodGet, path: "/api/test/test", want: "/api/test/:slug"},
		{name: "/test/:slug", method: http.MethodGet, path: "/test/test", want: "/test/:slug"},
		{name: "/test", method: http.MethodGet, path: "/test", want: "Get: /test"},
		{name: "/test", method: http.MethodPost, path: "/test", want: "Post: /test"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			assert.NoError(t, err)
			got := httptest.NewRecorder()
			router.Serve(got, req)
			assert.Equal(t, tt.want, got.Body.String())
		})
	}
}
