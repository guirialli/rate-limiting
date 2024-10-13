package entity

import "time"

type RaterLimit struct {
	Trys          int        `json:"trys"`
	Type          string     `json:"type"`
	AccessTimeout time.Time  `json:"access_timeout"`
	BlockAt       *time.Time `json:"block_at"`
}

func NewRaterLimit(typer string, timeout time.Duration) RaterLimit {
	return RaterLimit{
		Trys:          0,
		Type:          typer,
		AccessTimeout: time.Now().Add(timeout),
		BlockAt:       nil,
	}
}
