package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/internals/infra/webserver/controller"
)

type Book struct {
	Controller controller.IBooks
	auth       IAuthToken
}

func NewBook(controller controller.IBooks, auth IAuthToken) *Book {
	return &Book{
		Controller: controller,
		auth:       auth,
	}
}

func (b *Book) Use(r *chi.Mux) error {

	r.Route("/books", func(r chi.Router) {
		r.Use(jwtauth.Verifier(b.auth.NewTokenAuth()))
		r.Use(jwtauth.Authenticator)

		r.Get("/", b.Controller.GetAll)
		r.Get("/author", b.Controller.GetAllWithAuthor)
		r.Get("/{id}", b.Controller.GetById)
		r.Get("/{id}/author", b.Controller.GetByIdWithAuthor)

		r.Post("/", b.Controller.Create)
		r.Put("/{id}", b.Controller.Update)
		r.Patch("/{id}", b.Controller.Patch)
		r.Delete("/{id}", b.Controller.Delete)
	})

	r.Route("/public/books", func(r chi.Router) {
		r.Get("/", b.Controller.GetAll)
		r.Get("/{id}", b.Controller.GetByIdWithAuthor)
	})
	return nil
}
