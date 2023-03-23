package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type RouterInterface interface {
	Handler() http.Handler
	HandleFunc(path string, method string, fn func(w http.ResponseWriter, r *http.Request))
}

type Listener interface {
	RegisterEndpoints(router RouterInterface)
}

type Router struct {
	router *mux.Router
}

func NewRouter() RouterInterface {
	return &Router{router: mux.NewRouter()}
}

func (r *Router) HandleFunc(path string, method string, f func(http.ResponseWriter, *http.Request)) {
	r.router.HandleFunc(path, f).Methods(method)
}

func (r *Router) Handler() http.Handler {
	r.router.Methods(http.MethodOptions).
		HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
		})

	return r.router
}
