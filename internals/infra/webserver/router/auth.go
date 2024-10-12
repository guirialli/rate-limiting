package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/guirialli/rater_limit/internals/infra/webserver/controller"
)

type Auth struct {
	Controller controller.IAuth
}

func NewAuth(controller controller.IAuth) *Auth {
	return &Auth{
		Controller: controller,
	}
}

func (a *Auth) Use(r *chi.Mux) error {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", a.Controller.Register)
		r.Post("/login", a.Controller.Login)
	})
	return nil
}
