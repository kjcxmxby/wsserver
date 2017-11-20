package router

import (
	"net/http"
	"wsserver/log"
)

func init() {
	default_router.router = make(map[string]handler_fun)
}

type handler_fun func(*http.Request, http.ResponseWriter)

var (
	default_router Router
)

type Router struct {
	router map[string]handler_fun
}

func (r *Router) Register(path string, f handler_fun) {
	r.router[path] = f
}

func (r *Router) Router(path string, req *http.Request, w http.ResponseWriter) {
	f, ok := r.router[path]

	if ok {
		go f(req, w)
	}

	log.Debug("router path :", path)
}

func Register(path string, f handler_fun) {
	default_router.Register(path, f)
}

func RouterReq(path string, req *http.Request, w http.ResponseWriter) {
	default_router.Router(path, req, w)
}
