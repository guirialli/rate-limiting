package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Swagger struct {
	path string
}

func NewSwagger(path string) *Swagger {
	return &Swagger{path: path}
}

func (s *Swagger) Use(r *chi.Mux) error {
	path := fmt.Sprintf("%s/swagger.json", s.path)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
	return nil
}
