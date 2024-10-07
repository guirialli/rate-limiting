package dtos

type BookCreate struct {
	Title       string
	Pages       int
	Description *string
	Author      string
}

type BookUpdate struct {
	Title       string
	Pages       int
	Description *string
	Author      string
}

type BookPatch struct {
	Title       *string
	Pages       *int
	Description *string
	Author      *string
}
