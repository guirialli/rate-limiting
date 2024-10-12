package dtos

import "github.com/guirialli/rater_limit/internals/entity"

type AuthorWithBooks struct {
	Author ResponseAuthor `json:"author"`
	Books  []entity.Book  `json:"books"`
}

type BookWithAuthor struct {
	Book   entity.Book    `json:"book"`
	Author ResponseAuthor `json:"author"`
}
