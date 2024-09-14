package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Handlers   http.HandlerFunc
	HttpMethod string
	Path       string
}

type WithMiddlewares struct {
	Routers     []Router
	Middlewares []func(next http.Handler) http.Handler
}

type ChiRoute struct {
	Path                  string
	Route                 chi.Route
	RouterWithMiddlewares WithMiddlewares
}
