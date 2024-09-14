package entity

import "github.com/google/uuid"

type Book struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Pages       int     `json:"pages"`
	Description *string `json:"description"`
	Author      string  `json:"author"`
}

func NewBook(title string, pages int, author string, description *string) *Book {
	return &Book{
		Id:          uuid.NewString(),
		Title:       title,
		Pages:       pages,
		Description: description,
		Author:      author,
	}
}
