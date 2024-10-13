package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
)

// Auth handles authentication requests.
type Auth struct {
	userUseCase usecases.IUser
	db          *sql.DB
	errHandler  IHttpHandlerError
}

// NewAuth creates a new Auth controller.
func NewAuth(db *sql.DB, userUseCase usecases.IUser, errHandler IHttpHandlerError) *Auth {
	return &Auth{
		userUseCase: userUseCase,
		db:          db,
		errHandler:  errHandler,
	}
}

// Login godoc
// @Summary Login user
// @Description Login to get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dtos.LoginForm true "Login Form"
// @Success 200 {object} dtos.ResponseJwt "Successful login"
// @Failure 400 {object} ErrorResponse "Invalid body"
// @Failure 401 {object} ErrorResponse "Invalid username or password"
// @Router /auth/login [post]
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

// Register godoc
// @Summary Register user
// @Description Register to create a new user and get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param register body dtos.RegisterForm true "Register Form"
// @Success 201 {object} dtos.ResponseJwt "Successful registration"
// @Failure 400 {object} ErrorResponse "Invalid body"
// @Router /auth/register [post]
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

// ErrorResponse is used for returning error messages.
type ErrorResponse struct {
	Message string `json:"message"`
}
