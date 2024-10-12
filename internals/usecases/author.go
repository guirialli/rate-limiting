package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"github.com/guirialli/rater_limit/pkg/uow"
)

type Author struct{}

func NewAuthor() *Author {
	return &Author{}
}

func (a *Author) scan(r *sql.Rows) (entity.Author, error) {
	var author entity.Author
	err := r.Scan(&author.Id, &author.Name, &author.Birthday, &author.Description)
	if err != nil {
		return author, fmt.Errorf("error scanning author: %w", err)
	}
	return author, err
}

func (a *Author) scanRows(rows *sql.Rows) ([]entity.Author, error) {
	authors := make([]entity.Author, 0)
	for rows.Next() {
		author, err := a.scan(rows)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (a *Author) FindAll(ctx context.Context, db *sql.DB) ([]entity.Author, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return a.scanRows(rows)
}

func (a *Author) FindAllWithBooks(ctx context.Context, db *sql.DB, bookUseCase IBook) ([]dtos.AuthorWithBooks, error) {
	authors, err := a.FindAll(ctx, db)
	if err != nil {
		return nil, err
	}
	authorsBook := make([]dtos.AuthorWithBooks, len(authors))
	for i, author := range authors {
		books, err := bookUseCase.FindAllByAuthor(ctx, db, author.Id)
		if err != nil {
			return nil, err
		}
		authorsBook[i].Author = *dtos.ConvertAuthorToAuthorResponse(&author)
		authorsBook[i].Books = books
	}
	return authorsBook, nil
}

func (a *Author) FindById(ctx context.Context, db *sql.DB, id string) (*entity.Author, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM authors WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, fmt.Errorf("error finding author by id: %s", id)
	}
	defer rows.Close()

	rows.Next()
	author, err := a.scan(rows)
	if err != nil {
		return nil, fmt.Errorf("error finding author by id: %s", id)
	}
	return &author, nil
}

func (a *Author) FindByIdWithBooks(ctx context.Context, db *sql.DB, id string, bookUseCase IBook) (*dtos.AuthorWithBooks, error) {
	author, err := a.FindById(ctx, db, id)
	if err != nil {
		return nil, fmt.Errorf("error finding author by id: %s", id)
	}

	books, err := bookUseCase.FindAllByAuthor(ctx, db, author.Id)
	if err != nil {
		return nil, err
	}

	var authorBooks dtos.AuthorWithBooks
	authorBooks.Author = *dtos.ConvertAuthorToAuthorResponse(author)
	authorBooks.Books = books

	return &authorBooks, err
}

func (a *Author) Create(ctx context.Context, db *sql.DB, authorCreate *dtos.AuthorCreate) (*entity.Author, error) {
	author := &entity.Author{
		Id:          uuid.NewString(),
		Name:        authorCreate.Name,
		Description: authorCreate.Description,
		Birthday:    authorCreate.Birthday,
	}

	_, err := uow.NewTransaction(db, func() (*entity.Author, error) {
		_, err := db.ExecContext(ctx, "INSERT INTO authors(id, name, description, birthday) VALUES (?, ?, ?, ?)",
			author.Id, author.Name, author.Description, author.Birthday)
		return author, err
	}).Exec()
	if err != nil {
		return nil, fmt.Errorf("error creating author: %w", err)
	}

	return author, nil
}

func (a *Author) Update(ctx context.Context, db *sql.DB, id string, authorUpdate *dtos.AuthorUpdate) (*entity.Author, error) {
	return uow.NewTransaction(db, func() (*entity.Author, error) {
		_, err := db.ExecContext(ctx, "UPDATE authors SET name=?, description=?, birthday=? WHERE id=?",
			authorUpdate.Name, authorUpdate.Description, authorUpdate.Birthday, id,
		)

		if err != nil {
			return nil, fmt.Errorf("error updating author: %w", err)
		}
		return a.FindById(ctx, db, id)
	}).Exec()
}

func (a *Author) Patch(ctx context.Context, db *sql.DB, id string, authorPatch *dtos.AuthorPatch) (*entity.Author, error) {
	author, err := a.FindById(ctx, db, id)
	if err != nil {
		return nil, err
	}

	if authorPatch.Name == nil {
		authorPatch.Name = &author.Name
	}
	if authorPatch.Description == nil {
		authorPatch.Description = author.Description
	}
	if authorPatch.Birthday == nil {
		authorPatch.Birthday = author.Birthday
	}

	authorUpdate := &dtos.AuthorUpdate{
		Name:        *authorPatch.Name,
		Birthday:    authorPatch.Birthday,
		Description: authorPatch.Description,
	}

	return a.Update(ctx, db, id, authorUpdate)
}

func (a *Author) Delete(ctx context.Context, db *sql.DB, id string) error {
	if _, err := a.FindById(ctx, db, id); err != nil {
		return err
	}

	_, err := uow.NewTransaction(db, func() (interface{}, error) {
		_, err := db.ExecContext(ctx, "DELETE FROM authors WHERE id=?", id)
		return nil, err
	}).Exec()
	return err
}
