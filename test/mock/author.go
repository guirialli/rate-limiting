package mock

import (
	"github.com/guirialli/rater_limit/internals/vos"
	"time"
)

type Author struct{}

func NewAuthor() *Author {
	return &Author{}
}

func (a *Author) Create() *vos.AuthorCreate {
	birthday := time.Now()
	return &vos.AuthorCreate{
		Name:        "Test Author",
		Description: nil,
		Birthday:    &birthday,
	}
}
