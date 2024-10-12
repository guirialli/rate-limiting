package dtos

import "time"

type AuthorCreate struct {
	Name        string     `json:"name"`
	Birthday    *time.Time `json:"birthday"`
	Description *string    `json:"description"`
}

type AuthorUpdate struct {
	Name        string     `json:"name"`
	Birthday    *time.Time `json:"birthday"`
	Description *string    `json:"description"`
}

type AuthorPatch struct {
	Name        *string    `json:"name"`
	Birthday    *time.Time `json:"birthday"`
	Description *string    `json:"description"`
}

type AuthorBody struct {
	Name        *string `json:"name"`
	Birthday    *string `json:"birthday"`
	Description *string `json:"description"`
}
