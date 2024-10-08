package dtos

type ResponseJson[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

type ResponseJwt struct {
	Token string `json:"token"`
}
