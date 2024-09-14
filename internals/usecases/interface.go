package usecases

import (
	"context"
	"database/sql"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/vos"
)

type IBook interface {
	FindAll(ctx context.Context, db sql.DB) (*[]entity.Book, error)
	FindAllWithAuthor(ctx context.Context, db sql.DB) (*[]vos.BookWithAuthor, error)
	Create(ctx context.Context, db sql.DB, bookCreate *vos.BookCreate) (*entity.Book, error)
	Patch(ctx context.Context, db sql.DB, id string, bookUpdate *vos.BookPatch) (*entity.Book, error)
	Update(ctx context.Context, db sql.DB, id string, bookUpdate *vos.BookUpdate) (*entity.Book, error)
	Delete(ctx context.Context, db sql.DB, id string) error
}

type IAuthor interface {
	FindAll(ctx context.Context, db sql.DB) (*[]entity.Author, error)
	FindAllWithBooks(ctx context.Context, db sql.DB) (*vos.AuthorWithBooks, error)
	Create(ctx context.Context, db sql.DB, authorCreate *vos.AuthorCreate) (*entity.Author, error)
	Patch(ctx context.Context, db sql.DB, id string, authorUpdate *vos.AuthorPatch) (*entity.Author, error)
	Update(ctx context.Context, db sql.DB, id string, authorUpdate *vos.AuthorUpdate) (*entity.Author, error)
	Delete(ctx context.Context, db sql.DB, id string) error
}
