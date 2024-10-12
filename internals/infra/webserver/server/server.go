package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/guirialli/rater_limit/config"
	"github.com/guirialli/rater_limit/internals/infra/webserver/middleware"
	"github.com/guirialli/rater_limit/internals/infra/webserver/router"
	"net/http"
)

type Server struct {
	config  *config.WebServer
	routers []router.UseRouter
}

func NewServer(config *config.WebServer, routers []router.UseRouter) *Server {
	return &Server{
		config:  config,
		routers: routers,
	}
}

func (s *Server) Start() error {
	cfgRater, err := config.LoadRaterLimitConfig()
	if err != nil {
		return err
	}
	rateLimit, err := middleware.NewRaterLimit(*cfgRater)
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Use(rateLimit.RateLimitMiddleware())

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	for _, route := range s.routers {
		if err = route.Use(r); err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s:%d", s.config.Ip, s.config.Port)
	fmt.Printf("Server listen on http://%s\n", url)
	if err = http.ListenAndServe(url, r); err != nil {
		return err
	}
	return nil
}
