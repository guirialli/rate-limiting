package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
	"time"
)

type Author struct {
	db          *sql.DB
	useCase     usecases.IAuthor
	bookUseCase usecases.IBook
	errHandler  IHttpHandlerError
}

func NewAuthor(db *sql.DB, useCase usecases.IAuthor, book usecases.IBook, errHandler IHttpHandlerError) *Author {
	return &Author{
		db:          db,
		useCase:     useCase,
		bookUseCase: book,
		errHandler:  errHandler,
	}
}

func (a *Author) unmarshalBody(r *http.Request) (*dtos.AuthorPatch, error) {
	var body dtos.AuthorBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	var birthday *time.Time = nil
	if body.Birthday != nil {
		dt, err := time.Parse("2006-01-02", *body.Birthday)
		if err != nil {
			return nil, err
		}
		birthday = &dt
	}
	return &dtos.AuthorPatch{
		Name:        body.Name,
		Birthday:    birthday,
		Description: body.Description,
	}, nil

}

func (a *Author) response(w http.ResponseWriter, author any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(dtos.ResponseJson[any]{
		Status: statusCode,
		Data:   author,
	})

	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Author) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := a.useCase.FindAll(ctx, a.db)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseAuthors := make([]dtos.ResponseAuthor, len(authors))
	for i, author := range authors {
		responseAuthors[i] = *dtos.ConvertAuthorToAuthorResponse(&author)
	}

	a.response(w, responseAuthors, http.StatusOK)
}

func (a *Author) GetAllWithBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := a.useCase.FindAllWithBooks(ctx, a.db, a.bookUseCase)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.response(w, authors, http.StatusOK)

}

func (a *Author) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	author, err := a.useCase.FindById(r.Context(), a.db, id)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	a.response(w, dtos.ConvertAuthorToAuthorResponse(author), http.StatusOK)
}

func (a *Author) GetByIdWithBooks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	author, err := a.useCase.FindByIdWithBooks(r.Context(), a.db, id, a.bookUseCase)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dtos.ResponseJson[dtos.AuthorWithBooks]{
		Status: http.StatusOK,
		Data:   *author,
	})

	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Author) Create(w http.ResponseWriter, r *http.Request) {
	body, err := a.unmarshalBody(r)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	} else if body.Name == nil {
		a.errHandler.ResponseError(w, "name is required", http.StatusBadRequest)
		return
	}

	author, err := a.useCase.Create(r.Context(), a.db, &dtos.AuthorCreate{
		Name:        *body.Name,
		Birthday:    body.Birthday,
		Description: body.Description,
	})

	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.response(w, dtos.ConvertAuthorToAuthorResponse(author), http.StatusCreated)
}

func (a *Author) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	body, err := a.unmarshalBody(r)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, err = a.useCase.FindById(r.Context(), a.db, id); err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	} else if body.Name == nil {
		a.errHandler.ResponseError(w, "name is required", http.StatusBadRequest)
		return
	}

	author, err := a.useCase.Update(r.Context(), a.db, id, &dtos.AuthorUpdate{
		Name:        *body.Name,
		Birthday:    body.Birthday,
		Description: body.Description,
	})
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.response(w, dtos.ConvertAuthorToAuthorResponse(author), http.StatusOK)
}

func (a *Author) Patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	body, err := a.unmarshalBody(r)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	} else if _, err = a.useCase.FindById(r.Context(), a.db, id); err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	author, err := a.useCase.Patch(r.Context(), a.db, id, body)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.response(w, dtos.ConvertAuthorToAuthorResponse(author), http.StatusOK)
}

func (a *Author) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		a.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	err := a.useCase.Delete(r.Context(), a.db, id)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode([]byte{})

	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}
