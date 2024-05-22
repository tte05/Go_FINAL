package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
)

func NewLimiter(rateLimit rate.Limit, burstLimit int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rateLimit, burstLimit)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() == false {
				http.Error(w, "Too Many Requests. Try again later.", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
