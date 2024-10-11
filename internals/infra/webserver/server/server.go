package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/webserver/middleware"
	"github.com/guirialli/rater_limit/internals/infra/webserver/router"
	"net/http"
	"strings"
)

type Server struct {
	routers     []router.Router
	middlewares []middleware.Middleware
	config      config.WebServer
	chiRoutes   []router.ChiRoute
}

func NewServer(
	router []router.Router,
	middlewares []middleware.Middleware,
	config config.WebServer, chiRoutes []router.ChiRoute) *Server {
	return &Server{
		routers:     router,
		middlewares: middlewares,
		config:      config,
		chiRoutes:   chiRoutes,
	}
}

func (s *Server) useRouters(r *chi.Mux, routers []router.Router) error {
	for _, ro := range routers {
		switch strings.ToUpper(ro.HttpMethod) {
		case "GET":
			r.Get(ro.Path, ro.Handlers)
		case "POST":
			r.Post(ro.Path, ro.Handlers)
		case "PUT":
			r.Put(ro.Path, ro.Handlers)
		case "PATCH":
			r.Patch(ro.Path, ro.Handlers)
		case "DELETE":
			r.Delete(ro.Path, ro.Handlers)
		default:
			return fmt.Errorf("invalid http method: %s", ro.HttpMethod)
		}

	}
	return nil
}

func (s *Server) useMiddlewares(r *chi.Mux, middlewares []middleware.Middleware) {
	for _, mw := range middlewares {
		r.Use(mw)
	}
}

func (s *Server) useChiRoute(r *chi.Mux, chiRoutes []router.ChiRoute) {
	for _, chiRo := range chiRoutes {
		r.Route(chiRo.Path, func(r chi.Router) {
			for _, mw := range chiRo.RouterWithMiddlewares.Middlewares {
				r.Use(mw)
			}
			for _, ro := range chiRo.RouterWithMiddlewares.Routers {
				switch strings.ToUpper(ro.HttpMethod) {
				case "GET":
					r.Get(ro.Path, ro.Handlers)
				case "POST":
					r.Post(ro.Path, ro.Handlers)
				case "PUT":
					r.Put(ro.Path, ro.Handlers)
				case "PATCH":
					r.Patch(ro.Path, ro.Handlers)
				case "DELETE":
					r.Delete(ro.Path, ro.Handlers)
				default:
					panic(fmt.Errorf("invalid http method: %s", ro.HttpMethod))
				}
			}
		})
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	if err := s.useRouters(r, s.routers); err != nil {
		return err
	}
	s.useMiddlewares(r, s.middlewares)
	s.useChiRoute(r, s.chiRoutes)

	url := fmt.Sprintf("http://%s:%d", s.config.Ip, s.config.Port)
	if err := http.ListenAndServe(url, r); err != nil {
		return err
	}
	return nil
}
