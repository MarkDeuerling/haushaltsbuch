package middleware

import "net/http"

// RecoverPanic ...
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// logging
			}
		}()
		next.ServeHTTP(w, r)
	})
}
