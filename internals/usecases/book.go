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

type Book struct {
}

func NewBook() *Book {
	return &Book{}
}

func (b *Book) scan(row *sql.Rows) (entity.Book, error) {
	var book entity.Book
	err := row.Scan(&book.Id, &book.Title, &book.Pages, &book.Description, &book.Author)
	if err != nil {
		return book, fmt.Errorf("error on scan data: %w", err)
	}
	return book, err
}

func (b *Book) scanRows(rows *sql.Rows) ([]entity.Book, error) {
	books := make([]entity.Book, 0)
	for rows.Next() {
		book, err := b.scan(rows)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (b *Book) FindAll(ctx context.Context, db *sql.DB) ([]entity.Book, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return b.scanRows(rows)
}

func (b *Book) FindAllWithAuthor(ctx context.Context, db *sql.DB, authorUseCases IAuthor) ([]dtos.BookWithAuthor, error) {
	var bookAuthors []dtos.BookWithAuthor
	books, err := b.FindAll(ctx, db)
	if err != nil {
		return nil, err
	}
	for _, book := range books {
		author, err := authorUseCases.FindById(ctx, db, book.Author)
		if err != nil {
			return nil, err
		}
		bookAuthors = append(bookAuthors, dtos.BookWithAuthor{Author: *author, Book: book})
	}
	return bookAuthors, err
}

func (b *Book) FindAllByAuthor(ctx context.Context, db *sql.DB, author string) ([]entity.Book, error) {
	rows, err := db.QueryContext(ctx, "SELECT * FROM books WHERE author_id=?", author)
	if err != nil {
		return nil, fmt.Errorf("failed to find book by author: %w", err)
	}
	defer rows.Close()
	return b.scanRows(rows)
}

func (b *Book) FindById(ctx context.Context, db *sql.DB, id string) (*entity.Book, error) {
	row, err := db.QueryContext(ctx, "SELECT * FROM books WHERE id=? LIMIT 1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to find book by id: %w", err)
	}
	defer row.Close()

	row.Next()
	book, err := b.scan(row)
	return &book, err
}

func (b *Book) FindByIdWithAuthor(ctx context.Context, db *sql.DB, id string, authorUseCases IAuthor) (*dtos.BookWithAuthor, error) {
	book, err := b.FindById(ctx, db, id)
	if err != nil {
		return nil, err
	}
	author, err := authorUseCases.FindById(ctx, db, book.Author)
	if err != nil {
		return nil, err
	}

	return &dtos.BookWithAuthor{Book: *book, Author: *author}, nil
}

func (b *Book) Create(ctx context.Context, db *sql.DB, bookCreate *dtos.BookCreate) (*entity.Book, error) {
	book := &entity.Book{
		Id:          uuid.New().String(),
		Title:       bookCreate.Title,
		Description: bookCreate.Description,
		Pages:       bookCreate.Pages,
		Author:      bookCreate.Author,
	}

	transaction := uow.NewTransaction(db, func() (*entity.Book, error) {
		_, err := db.ExecContext(ctx,
			"INSERT INTO books(id, title, pages, description, author_id) VALUES (?, ?, ?, ?, ?)",
			book.Id, book.Title, book.Pages, book.Description, book.Author)
		return book, err
	})
	_, err := transaction.Exec()
	if err != nil {
		return nil, fmt.Errorf("error executing insert book %v", err)
	}

	return book, nil
}

func (b *Book) Update(ctx context.Context, db *sql.DB, id string, bookUpdate *dtos.BookUpdate) (*entity.Book, error) {
	transaction := uow.NewTransaction(db, func() (*entity.Book, error) {
		_, err := db.ExecContext(ctx, `UPDATE books SET 
			title = ?, pages =?, description = ? ,author_id= ? WHERE id=?`,
			bookUpdate.Title, bookUpdate.Pages, bookUpdate.Description, bookUpdate.Author, id)
		if err != nil {
			return nil, fmt.Errorf("error on updating book: %w", err)
		}

		book, err := b.FindById(ctx, db, id)
		return book, err
	})
	return transaction.Exec()
}

func (b *Book) Patch(ctx context.Context, db *sql.DB, id string, bookPatch *dtos.BookPatch) (*entity.Book, error) {
	book, err := b.FindById(ctx, db, id)
	if err != nil {
		return nil, err
	}
	var bookUpdate dtos.BookUpdate

	if bookPatch.Title != nil {
		bookUpdate.Title = *bookPatch.Title
	} else {
		bookUpdate.Title = book.Title
	}

	if bookPatch.Pages != nil {
		bookUpdate.Pages = *bookPatch.Pages
	} else {
		bookUpdate.Pages = book.Pages
	}

	if bookPatch.Description != nil {
		bookUpdate.Description = bookPatch.Description
	} else {
		bookUpdate.Description = book.Description
	}

	if bookPatch.Author != nil {
		bookUpdate.Author = *bookPatch.Author
	} else {
		bookUpdate.Author = book.Author
	}

	return b.Update(ctx, db, id, &bookUpdate)
}

func (b *Book) Delete(ctx context.Context, db *sql.DB, id string) error {
	transaction := uow.NewTransaction(db, func() (*entity.Book, error) {
		if _, err := db.ExecContext(ctx, "DELETE FROM books where id=?", id); err != nil {
			return nil, fmt.Errorf("error on delete book: %w", err)
		}
		return nil, nil
	})

	_, err := transaction.Exec()
	return err

}
