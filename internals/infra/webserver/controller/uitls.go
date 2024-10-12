package controller

import (
	"encoding/json"
	"github.com/guirialli/rater_limit/internals/entity/dtos"
	"net/http"
)

type Utils struct {
}

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) ResponseError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(dtos.ResponseError{"error": message, "status": statusCode})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
