package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
	"net/http/httptest"
)

type Auth struct {
	userUseCase *usecases.User
	db          *sql.DB
}

func NewAuth(userUseCase *usecases.User, db *sql.DB) *Auth {
	return &Auth{
		userUseCase: userUseCase,
		db:          db,
	}
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var formLogin dtos.LoginForm
	err := json.NewDecoder(r.Body).Decode(&formLogin)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	token, err := a.userUseCase.Login(context.Background(), a.db, &formLogin)
	if err != nil {
		http.Error(w, "username or password invalid", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[dtos.ResponseJwt]{
		Status: http.StatusOK,
		Data: dtos.ResponseJwt{
			Token: token,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Auth) Register(w *httptest.ResponseRecorder, r *http.Request) {
	var formRegister dtos.RegisterForm
	err := json.NewDecoder(r.Body).Decode(&formRegister)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	token, err := a.userUseCase.Register(context.Background(), a.db, &formRegister)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[dtos.ResponseJwt]{
		Status: http.StatusCreated,
		Data: dtos.ResponseJwt{
			Token: token,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
