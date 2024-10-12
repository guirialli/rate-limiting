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

	r.Route("/authors", func(r chi.Router) {
		r.Use(jwtauth.Verifier(a.auth.NewTokenAuth()))
		r.Use(jwtauth.Authenticator)

		r.Get("/", a.Controller.GetAll)
		r.Get("/books", a.Controller.GetAllWithBooks)
		r.Get("/{id}", a.Controller.GetById)
		r.Get("/{id}/books", a.Controller.GetByIdWithBooks)

		r.Post("/", a.Controller.Create)
		r.Put("/{id}", a.Controller.Update)
		r.Patch("/{id}", a.Controller.Patch)
		r.Delete("/{id}", a.Controller.Delete)
	})

	r.Route("/public/authors", func(r chi.Router) {
		r.Get("/", a.Controller.GetAll)
		r.Get("/{id}", a.Controller.GetByIdWithBooks)
	})
	return nil
}
