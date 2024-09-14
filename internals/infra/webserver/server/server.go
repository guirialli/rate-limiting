package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/webserver/router"
	"net/http"
	"strings"
)

type Server struct {
	routers     []router.Router
	middlewares []func(next http.Handler) http.Handler
	config      config.WebServer
	chiRoutes   []router.ChiRoute
}

func NewServer(
	router []router.Router,
	middlewares []func(next http.Handler) http.Handler,
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
		if err := s.validHttpMethod(ro.HttpMethod); err != nil {
			return err
		}
		method := strings.ToUpper(ro.HttpMethod)
		if method == "GET" {
			r.Get(ro.Path, ro.Handlers)
		} else if method == "POST" {
			r.Post(ro.Path, ro.Handlers)
		} else if method == "PUT" {
			r.Put(ro.Path, ro.Handlers)
		} else if method == "PATCH" {
			r.Patch(ro.Path, ro.Handlers)
		} else if method == "DELETE" {
			r.Delete(ro.Path, ro.Handlers)
		}
	}
	return nil
}

func (s *Server) useMiddlewares(r *chi.Mux, middlewares []func(next http.Handler) http.Handler) {
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
				method := strings.ToUpper(ro.HttpMethod)
				if err := s.validHttpMethod(method); err != nil {
					fmt.Printf("Fail on add router %s in %s\n", chiRo.Path, ro.Path)
					panic(err)
				}

				if method == "GET" {
					r.Get(ro.Path, ro.Handlers)
				} else if method == "POST" {
					r.Post(ro.Path, ro.Handlers)
				} else if method == "PUT" {
					r.Put(ro.Path, ro.Handlers)
				} else if method == "PATCH" {
					r.Patch(ro.Path, ro.Handlers)
				} else if method == "DELETE" {
					r.Delete(ro.Path, ro.Handlers)
				}
			}
		})
	}
}

func (s *Server) validHttpMethod(method string) error {
	methodUpper := strings.ToUpper(method)
	if methodUpper == "GET" || methodUpper == "POST" ||
		methodUpper == "PUT" || methodUpper == "PATCH" ||
		methodUpper == "DELETE" {
		return nil
	}
	return fmt.Errorf("invalid http method: %s", method)
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
