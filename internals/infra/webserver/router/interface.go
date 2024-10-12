package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

type UseRouter interface {
	Use(mux *chi.Mux) error
}

type IAuthToken interface {
	NewTokenAuth() *jwtauth.JWTAuth
}
