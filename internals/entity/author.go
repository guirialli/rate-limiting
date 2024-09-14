package entity

import (
	"github.com/google/uuid"
	"time"
)

type Author struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Birthday    *time.Time `json:"birthday"`
	Description *string    `json:"description"`
}

func NewAuthor(name string, birthday *time.Time, description *string) *Author {
	return &Author{
		Id:          uuid.NewString(),
		Name:        name,
		Birthday:    birthday,
		Description: description,
	}
}
