package mock

import (
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"time"
)

type Author struct{}

func NewAuthor() *Author {
	return &Author{}
}

func (a *Author) Create(description *string) *dtos.AuthorCreate {
	birthday := time.Now()
	return &dtos.AuthorCreate{
		Name:        "Test Author",
		Description: description,
		Birthday:    &birthday,
	}
}

func (a *Author) CreateBody(description *string) *dtos.AuthorBody {
	birthday := "2022-01-01"
	name := "Test Author"
	return &dtos.AuthorBody{
		Name:        &name,
		Description: description,
		Birthday:    &birthday,
	}
}

func (a *Author) Patch(name, description *string, birthday *time.Time) *dtos.AuthorPatch {
	if name == nil {
		n := "test patch"
		name = &n
	}
	if birthday == nil {
		b := time.Now()
		birthday = &b
	}
	if description == nil {
		d := "test patch lorem ipsum"
		description = &d
	}

	return &dtos.AuthorPatch{
		Name:        name,
		Birthday:    birthday,
		Description: description,
	}
}
