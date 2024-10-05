package mock

import (
	"github.com/guirialli/rater_limit/internals/vos"
	"time"
)

type Author struct{}

func NewAuthor() *Author {
	return &Author{}
}

func (a *Author) Create(description *string) *vos.AuthorCreate {
	birthday := time.Now()
	return &vos.AuthorCreate{
		Name:        "Test Author",
		Description: description,
		Birthday:    &birthday,
	}
}

func (a *Author) Patch(name, description *string, birthday *time.Time) *vos.AuthorPatch {
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

	return &vos.AuthorPatch{
		Name:        name,
		Birthday:    birthday,
		Description: description,
	}
}
