package middleware

import (
	"github.com/guirialli/rater_limit/internals/usecases"
	"net/http"
	"strings"
)

type RaterLimit struct {
	useCase usecases.IRaterLimit
}

func NewRaterLimit(useCase usecases.IRaterLimit) *RaterLimit {
	return &RaterLimit{useCase}
}

func (rl *RaterLimit) RateLimitMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Authorization")
			if !rl.useCase.ValidToken(key) {
				key = strings.Split(r.RemoteAddr, ":")[0]
			}

			if !rl.useCase.TrackAccess(key) {
				http.Error(w,
					"you have reached the maximum number of requests or actions allowed within a certain time frame",
					http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
