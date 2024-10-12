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
	r.Post("/auth/login", a.Controller.Login)
	r.Post("/auth/register", a.Controller.Register)
	return nil
}
