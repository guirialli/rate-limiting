package dtos

import "github.com/guirialli/rater_limit/internals/entity"

func ConvertAuthorToAuthorResponse(a *entity.Author) *ResponseAuthor {
	response := ResponseAuthor{
		Id:          &a.Id,
		Name:        &a.Name,
		Birthday:    nil,
		Description: a.Description,
	}

	if a.Birthday != nil {
		dt := a.Birthday.Format("2006-01-02")
		response.Birthday = &dt
	}
	return &response
}
