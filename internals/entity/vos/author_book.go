package vos

import "github.com/guirialli/rater_limit/internals/entity"

type AuthorWithBooks struct {
	Author entity.Author `json:"author"`
	Books  []entity.Book `json:"books"`
}

type BookWithAuthor struct {
	Book   entity.Book   `json:"book"`
	Author entity.Author `json:"author"`
}
