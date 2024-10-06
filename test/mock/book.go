package mock

import (
	"github.com/guirialli/rater_limit/internals/entity/vos"
)

type Book struct{}

func NewBookMock() *Book {
	return &Book{}
}

func (b *Book) Create(description *string) *vos.BookCreate {
	return &vos.BookCreate{
		Title:       "test",
		Pages:       100,
		Description: description,
		Author:      "test",
	}
}

func (b *Book) CreateWithAuthor(author string, description *string) *vos.BookCreate {
	return &vos.BookCreate{
		Title:       "test",
		Pages:       100,
		Description: description,
		Author:      author,
	}
}
