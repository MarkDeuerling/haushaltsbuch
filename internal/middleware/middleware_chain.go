package middleware

import "net/http"

type middleware func(http.Handler) http.Handler

// Chain ...
func Chain(m ...middleware) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		next := handler
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
