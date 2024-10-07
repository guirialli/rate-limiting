package usecases

import (
	"context"
	"database/sql"
	"github.com/go-chi/jwtauth"
	"github.com/guirialli/rater_limit/internals/entity"
	vos "github.com/guirialli/rater_limit/internals/entity/dtos"
)

type IBook interface {
	FindAll(ctx context.Context, db *sql.DB) ([]entity.Book, error)
	FindAllWithAuthor(ctx context.Context, db *sql.DB, authorUseCase *Author) ([]vos.BookWithAuthor, error)
	FindAllByAuthor(ctx context.Context, db *sql.DB, author string) ([]entity.Book, error)
	FindById(ctx context.Context, db *sql.DB, id string) (*entity.Book, error)
	FindByIdWithAuthor(ctx context.Context, db *sql.DB, id string, authorUseCases *Author) (*vos.BookWithAuthor, error)
	Create(ctx context.Context, db *sql.DB, bookCreate *vos.BookCreate) (*entity.Book, error)
	Patch(ctx context.Context, db *sql.DB, id string, bookUpdate *vos.BookPatch) (*entity.Book, error)
	Update(ctx context.Context, db *sql.DB, id string, bookUpdate *vos.BookUpdate) (*entity.Book, error)
	Delete(ctx context.Context, db *sql.DB, id string) error
}

type IAuthor interface {
	FindAll(ctx context.Context, db *sql.DB) ([]entity.Author, error)
	FindAllWithBooks(ctx context.Context, db *sql.DB, bookUseCase *Book) ([]vos.AuthorWithBooks, error)
	FindById(ctx context.Context, db *sql.DB, id string) (*entity.Author, error)
	FindByIdWithBooks(ctx context.Context, db *sql.DB, id string, bookUseCase *Book) (*vos.AuthorWithBooks, error)
	Create(ctx context.Context, db *sql.DB, authorCreate *vos.AuthorCreate) (*entity.Author, error)
	Patch(ctx context.Context, db *sql.DB, id string, authorUpdate *vos.AuthorPatch) (*entity.Author, error)
	Update(ctx context.Context, db *sql.DB, id string, authorUpdate *vos.AuthorUpdate) (*entity.Author, error)
	Delete(ctx context.Context, db *sql.DB, id string) error
}

type IUser interface {
	Login(ctx context.Context, db *sql.DB, form *vos.LoginForm) (string, error)
	Register(ctx context.Context, db *sql.DB, form *vos.RegisterForm) (string, error)
	NewTokenAuth() *jwtauth.JWTAuth
}
