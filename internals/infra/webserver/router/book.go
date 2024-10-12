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
	r.Get("/books", b.Controller.GetAll)
	r.Get("/books/author", b.Controller.GetAllWithAuthor)
	r.Get("/books/{id}", b.Controller.GetById)
	r.Get("/books/{id}/author", b.Controller.GetByIdWithAuthor)
	r.Route("/books", func(r chi.Router) {
		r.Use(jwtauth.Verifier(b.auth.NewTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Post("/", b.Controller.Create)
		r.Put("/{id}", b.Controller.Update)
		r.Patch("/{id}", b.Controller.Update)
		r.Delete("/{id}", b.Controller.Delete)
	})
	return nil
}
