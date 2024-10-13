package middleware

import "net/http"

type IRaterLimit interface {
	RateLimitMiddleware() func(next http.Handler) http.Handler
}
