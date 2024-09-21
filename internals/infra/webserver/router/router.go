package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/guirialli/rater_limit/internals/infra/webserver/middleware"
	"net/http"
)

type Router struct {
	Handlers   http.HandlerFunc
	HttpMethod string
	Path       string
}

type WithMiddlewares struct {
	Routers     []Router
	Middlewares []middleware.Middleware
}

type ChiRoute struct {
	Path                  string
	Route                 chi.Route
	RouterWithMiddlewares WithMiddlewares
}
