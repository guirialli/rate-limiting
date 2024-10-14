package entity

import "time"

type RaterLimit struct {
	Id            string     `json:"id"`
	Trys          int        `json:"trys"`
	Typer         string     `json:"type"`
	AccessTimeout time.Time  `json:"access_timeout"`
	BlockAt       *time.Time `json:"block_at"`
}

func NewRaterLimit(typer string, timeout time.Duration) *RaterLimit {
	return &RaterLimit{
		Trys:          0,
		Typer:         typer,
		AccessTimeout: time.Now().Add(timeout),
		BlockAt:       nil,
	}
}

func NewRateLimitSql(id string, typer string, timeout time.Duration) *RaterLimit {
	return &RaterLimit{
		Id:            id,
		Trys:          0,
		Typer:         typer,
		AccessTimeout: time.Now().Add(timeout),
		BlockAt:       nil,
	}
}
