package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/internals/infra/webserver/controller"
)

type Author struct {
	Controller controller.IAuthor
	auth       IAuthToken
}

func NewAuthor(controller controller.IAuthor, auth IAuthToken) *Author {
	return &Author{
		Controller: controller,
		auth:       auth,
	}
}

func (a *Author) Use(r *chi.Mux) error {
	r.Get("/authors", a.Controller.GetAll)
	r.Get("/authors/books", a.Controller.GetAllWithBooks)
	r.Get("/authors/{id}", a.Controller.GetById)
	r.Get("/authors/{id}/books", a.Controller.GetByIdWithBooks)
	r.Route("/authors", func(r chi.Router) {
		r.Use(jwtauth.Verifier(a.auth.NewTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Post("/", a.Controller.Create)
		r.Put("/{id}", a.Controller.Update)
		r.Patch("/{id}", a.Controller.Update)
		r.Delete("/{id}", a.Controller.Delete)
	})
	return nil
}
