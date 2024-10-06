package usecases

import (
	"context"
	"database/sql"
	"github.com/guirialli/rater_limit/internals/entity"
	vos2 "github.com/guirialli/rater_limit/internals/entity/vos"
)

type IBook interface {
	FindAll(ctx context.Context, db *sql.DB) ([]entity.Book, error)
	FindById(ctx context.Context, db *sql.DB, id string) (*entity.Book, error)
	Create(ctx context.Context, db *sql.DB, bookCreate *vos2.BookCreate) (*entity.Book, error)
	Patch(ctx context.Context, db *sql.DB, id string, bookUpdate *vos2.BookPatch) (*entity.Book, error)
	Update(ctx context.Context, db *sql.DB, id string, bookUpdate *vos2.BookUpdate) (*entity.Book, error)
	Delete(ctx context.Context, db *sql.DB, id string) error
}

type IAuthor interface {
	FindAll(ctx context.Context, db *sql.DB) ([]entity.Author, error)
	FindById(ctx context.Context, db *sql.DB, id string) (*entity.Author, error)
	Create(ctx context.Context, db *sql.DB, authorCreate *vos2.AuthorCreate) (*entity.Author, error)
	Patch(ctx context.Context, db *sql.DB, id string, authorUpdate *vos2.AuthorPatch) (*entity.Author, error)
	Update(ctx context.Context, db *sql.DB, id string, authorUpdate *vos2.AuthorUpdate) (*entity.Author, error)
	Delete(ctx context.Context, db *sql.DB, id string) error
}

type IAuthorBook interface {
	FindAllWithAuthor(ctx context.Context, db *sql.DB) (*[]vos2.BookWithAuthor, error)
	FindAllWithBooks(ctx context.Context, db sql.DB) (*vos2.AuthorWithBooks, error)
}
