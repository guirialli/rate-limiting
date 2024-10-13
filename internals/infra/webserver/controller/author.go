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

// Author handles author-related requests.
type Author struct {
	db          *sql.DB
	useCase     usecases.IAuthor
	bookUseCase usecases.IBook
	errHandler  IHttpHandlerError
}

// NewAuthor creates a new Author controller.
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

// GetAll godoc
// @Summary Get all authors
// @Description Retrieve a list of all authors
// @Tags authors
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dtos.ResponseAuthor "List of authors"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /authors [get]
// @Router /public/authors [get]
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

// GetAllWithBooks godoc
// @Summary Get all authors with books
// @Description Retrieve a list of all authors along with their books
// @Tags authors
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dtos.AuthorWithBooks "List of authors with their books"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /authors/books [get]
func (a *Author) GetAllWithBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authors, err := a.useCase.FindAllWithBooks(ctx, a.db, a.bookUseCase)
	if err != nil {
		a.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.response(w, authors, http.StatusOK)
}

// GetById godoc
// @Summary Get author by ID
// @Description Retrieve a specific author by ID
// @Tags authors
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Success 200 {object} dtos.ResponseAuthor "Author details"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Author not found"
// @Router /authors/{id} [get]
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

// GetByIdWithBooks godoc
// @Summary Get author by ID with books
// @Description Retrieve a specific author by ID along with their books
// @Tags authors
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Success 200 {object} dtos.AuthorWithBooks "Author details with books"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Author not found"
// @Router /authors/{id}/books [get]
// @Router /public/authors/{id} [get]
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

	a.response(w, author, http.StatusOK)
}

// Create godoc
// @Summary Create a new author
// @Description Create a new author in the system
// @Tags authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param author body dtos.AuthorCreate true "Author data"
// @Success 201 {object} dtos.ResponseAuthor "Created author"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /authors [post]
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

// Update godoc
// @Summary Update an existing author
// @Description Update an existing author's details
// @Tags authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Param author body dtos.AuthorUpdate true "Author data"
// @Success 200 {object} dtos.ResponseAuthor "Updated author"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Author not found"
// @Router /authors/{id} [put]
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

// Patch godoc
// @Summary Partially update an author
// @Description Update specific fields of an author
// @Tags authors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Param author body dtos.AuthorPatch true "Author data"
// @Success 200 {object} dtos.ResponseAuthor "Updated author"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Author not found"
// @Router /authors/{id} [patch]
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

// Delete godoc
// @Summary Delete an author
// @Description Remove an author from the system
// @Tags authors
// @Security ApiKeyAuth
// @Param id path string true "Author ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Author not found"
// @Router /authors/{id} [delete]
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
