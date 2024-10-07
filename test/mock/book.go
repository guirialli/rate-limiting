package mock

import (
	"github.com/guirialli/rater_limit/internals/entity/dtos"
)

type Book struct{}

func NewBookMock() *Book {
	return &Book{}
}

func (b *Book) Create(description *string) *dtos.BookCreate {
	return &dtos.BookCreate{
		Title:       "test",
		Pages:       100,
		Description: description,
		Author:      "test",
	}
}

func (b *Book) CreateWithAuthor(author string, description *string) *dtos.BookCreate {
	return &dtos.BookCreate{
		Title:       "test",
		Pages:       100,
		Description: description,
		Author:      author,
	}
}
