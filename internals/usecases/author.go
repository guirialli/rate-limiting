package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/guirialli/rater_limit/internals/entity"
	"github.com/guirialli/rater_limit/internals/vos"
	"github.com/guirialli/rater_limit/pkg/uow"
)

type Author struct{}

func NewAuthor() *Author {
	return &Author{}
}

func (a *Author) scan(r *sql.Rows) (entity.Author, error) {
	var author entity.Author
	err := r.Scan(&author.Id, &author.Name, &author.Birthday, &author.Description)
	return author, err
}

func (a *Author) FindAll(ctx context.Context, db *sql.DB) ([]entity.Author, error) {
	var authors []entity.Author
	rows, err := db.QueryContext(ctx, "SELECT * FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		author, err := a.scan(rows)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}
func (a *Author) FindById(ctx context.Context, db *sql.DB, id string) (*entity.Author, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM authors WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	author, err := a.scan(rows)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (a *Author) Create(ctx context.Context, db *sql.DB, authorCreate *vos.AuthorCreate) (*entity.Author, error) {
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
