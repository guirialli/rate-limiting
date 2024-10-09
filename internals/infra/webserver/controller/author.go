package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/guirialli/rater_limit/internals/entity"
	vos "github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
)

type Author struct {
	db          *sql.DB
	useCase     usecases.IAuthor
	bookUseCase usecases.IBook
}

func NewAuthor(db *sql.DB, useCase usecases.IAuthor, book usecases.IBook) *Author {
	return &Author{
		db:          db,
		useCase:     useCase,
		bookUseCase: book,
	}
}

func (a *Author) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := a.useCase.FindAll(ctx, a.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(vos.ResponseJson[[]entity.Author]{
		Status: http.StatusOK,
		Data:   authors,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
