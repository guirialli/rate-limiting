package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
)

type Auth struct {
	userUseCase usecases.IUser
	db          *sql.DB
	errHandler  IHttpHandlerError
}

func NewAuth(db *sql.DB, userUseCase usecases.IUser, errHandler IHttpHandlerError) *Auth {
	return &Auth{
		userUseCase: userUseCase,
		db:          db,
		errHandler:  errHandler,
	}
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var formLogin dtos.LoginForm
	err := json.NewDecoder(r.Body).Decode(&formLogin)
	if err != nil {
		a.errHandler.ResponseError(w, "invalid body", http.StatusBadRequest)
		return
	}

	token, err := a.userUseCase.Login(context.Background(), a.db, &formLogin)
	if err != nil {
		a.errHandler.ResponseError(w, "invalid username or password", http.StatusUnauthorized)
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
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var formRegister dtos.RegisterForm
	err := json.NewDecoder(r.Body).Decode(&formRegister)
	if err != nil {
		a.errHandler.ResponseError(w, "invalid body", http.StatusBadRequest)
		return
	}

	token, err := a.userUseCase.Register(context.Background(), a.db, &formRegister)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
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
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
