package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
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
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]entity.Book]{
		Status: http.StatusOK,
		Data:   books,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Book) GetAllWithAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	books, err := b.useCase.FindAllWithAuthor(ctx, b.db, b.authorUseCase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]dtos.BookWithAuthor]{
		Status: http.StatusOK,
		Data:   books,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Book) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	book, err := b.useCase.FindById(r.Context(), b.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Book) Create(w http.ResponseWriter, r *http.Request) {
	var body dtos.BookCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := b.useCase.Create(r.Context(), b.db, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusCreated,
		Data:   *book,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Book) Update(w http.ResponseWriter, r *http.Request) {
	var body dtos.BookUpdate
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := b.useCase.Update(r.Context(), b.db, id, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Book) Path(w http.ResponseWriter, r *http.Request) {
	var body dtos.BookPatch
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := b.useCase.Patch(r.Context(), b.db, id, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (b *Book) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := b.useCase.Delete(r.Context(), b.db, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode([]byte{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
