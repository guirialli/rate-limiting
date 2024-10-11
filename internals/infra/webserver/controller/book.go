package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
)

type Book struct {
	db            *sql.DB
	useCase       usecases.IBook
	authorUseCase usecases.IAuthor
}

func NewBook(db *sql.DB, useCase usecases.IBook, authorUseCase usecases.IAuthor) *Book {
	return &Book{
		db:            db,
		useCase:       useCase,
		authorUseCase: authorUseCase,
	}
}

func (b *Book) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	books, err := b.useCase.FindAll(ctx, b.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]entity.Book]{
		Status: http.StatusOK,
		Data:   books,
	})
}
