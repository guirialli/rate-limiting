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

// BookController representa o controlador de livros
type Book struct {
	db            *sql.DB
	useCase       usecases.IBook
	authorUseCase usecases.IAuthor
	errHandler    IHttpHandlerError
}

// NewBook cria uma nova instância do controlador Book
func NewBook(db *sql.DB, useCase usecases.IBook, authorUseCase usecases.IAuthor, errHandler IHttpHandlerError) *Book {
	return &Book{
		db:            db,
		useCase:       useCase,
		authorUseCase: authorUseCase,
		errHandler:    errHandler,
	}
}

// GetAll godoc
// @Summary Get all books
// @Description Retrieve a list of all books
// @Tags books
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []entity.Book "List of books"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /books [get]
// @Router /public/books [get]
func (b *Book) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	ctx := r.Context()
	books, err := b.useCase.FindAll(ctx, b.db)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]entity.Book]{
		Status: http.StatusOK,
		Data:   books,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetById godoc
// @Summary Get book by ID
// @Description Retrieve a specific book by ID
// @Tags books
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Success 200 {object} entity.Book "Book details"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Router /books/{id} [get]
func (b *Book) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	book, err := b.useCase.FindById(r.Context(), b.db, id)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetByIdWithAuthos godoc
// @Summary Get book by ID with Author
// @Description Retrieve a specific book by ID
// @Tags books
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Success 200 {object} dtos.BookWithAuthor "Book details"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Router /books/{id}/author [get]
// @Router /public/books/{id} [get]
func (b *Book) GetByIdWithAuthor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	book, err := b.useCase.FindByIdWithAuthor(r.Context(), b.db, id, b.authorUseCase)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[dtos.BookWithAuthor]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAllWithAuthor godoc
// @Summary Get all books with authors
// @Description Retrieve a list of all books along with their authors
// @Tags books
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []dtos.BookWithAuthor "List of books with authors"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /books/author [get]
func (b *Book) GetAllWithAuthor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	books, err := b.useCase.FindAllWithAuthor(ctx, b.db, b.authorUseCase)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[[]dtos.BookWithAuthor]{
		Status: http.StatusOK,
		Data:   books,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create godoc
// @Summary Create a new book
// @Description Add a new book to the database
// @Tags books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param book body dtos.BookCreate true "Book data"
// @Success 201 {object} entity.Book "Created book"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Router /books [post]
func (b *Book) Create(w http.ResponseWriter, r *http.Request) {
	var body dtos.BookCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := b.useCase.Create(r.Context(), b.db, &body)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusCreated,
		Data:   *book,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// Update godoc
// @Summary Update an existing book
// @Description Update book details in the database
// @Tags books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Param book body dtos.BookUpdate true "Updated book data"
// @Success 200 {object} entity.Book "Updated book"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Router /books/{id} [put]
func (b *Book) Update(w http.ResponseWriter, r *http.Request) {
	var body dtos.BookUpdate
	id := chi.URLParam(r, "id")
	if id == "" {
		b.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := b.useCase.Update(r.Context(), b.db, id, &body)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// Patch godoc
// @Summary Patch an existing book
// @Description Partially update book details in the database
// @Tags books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Param book body dtos.BookPatch true "Patch data for the book"
// @Success 200 {object} entity.Book "Patched book"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Router /books/{id} [patch]
func (b *Book) Patch(w http.ResponseWriter, r *http.Request) {
	var body dtos.BookPatch
	id := chi.URLParam(r, "id")
	if id == "" {
		b.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := b.useCase.Patch(r.Context(), b.db, id, &body)
	if err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dtos.ResponseJson[entity.Book]{
		Status: http.StatusOK,
		Data:   *book,
	}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}

// Delete godoc
// @Summary Delete a book
// @Description Remove a book from the database
// @Tags books
// @Security ApiKeyAuth
// @Param id path string true "Book ID"
// @Success 204 "No content"
// @Failure 400 {object} ErrorResponse "ID is required"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Router /books/{id} [delete]
func (b *Book) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		b.errHandler.ResponseError(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := b.useCase.Delete(r.Context(), b.db, id); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode([]byte{}); err != nil {
		b.errHandler.ResponseError(w, err.Error(), http.StatusInternalServerError)
	}
}
