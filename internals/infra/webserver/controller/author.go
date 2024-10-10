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
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]entity.Author]{
		Status: http.StatusOK,
		Data:   authors,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Author) GetAllWithBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := a.useCase.FindAllWithBooks(ctx, a.db, a.bookUseCase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]dtos.AuthorWithBooks]{
		Status: http.StatusOK,
		Data:   authors,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Author) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	author, err := a.useCase.FindById(r.Context(), a.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Author]{
		Status: http.StatusOK,
		Data:   *author,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Author) GetByIdWithBooks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	author, err := a.useCase.FindByIdWithBooks(r.Context(), a.db, id, a.bookUseCase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[dtos.AuthorWithBooks]{
		Status: http.StatusOK,
		Data:   *author,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Author) Create(w http.ResponseWriter, r *http.Request) {
	var body dtos.AuthorCreate
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	author, err := a.useCase.Create(r.Context(), a.db, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Author]{
		Status: http.StatusCreated,
		Data:   *author,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Author) Update(w http.ResponseWriter, r *http.Request) {
	var body dtos.AuthorUpdate
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, err = a.useCase.FindById(r.Context(), a.db, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	author, err := a.useCase.Update(r.Context(), a.db, id, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Author]{
		Status: http.StatusOK,
		Data:   *author,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Author) Patch(w http.ResponseWriter, r *http.Request) {
	var body dtos.AuthorPatch
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, err = a.useCase.FindById(r.Context(), a.db, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	author, err := a.useCase.Patch(r.Context(), a.db, id, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Author]{
		Status: http.StatusOK,
		Data:   *author,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Author) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := a.useCase.Delete(r.Context(), a.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode([]byte{})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
