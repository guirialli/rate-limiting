package dtos

type ResponseJson[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

type ResponseJwt struct {
	Token string `json:"token"`
}

type ResponseAuthor struct {
	Id          *string `json:"id"`
	Name        *string `json:"name"`
	Birthday    *string `json:"birthday"`
	Description *string `json:"description"`
}

type ResponseError map[string]any
