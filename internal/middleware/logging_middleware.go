package middleware

import (
	"fmt"
	"net/http"

	"gitlab.com/shingeki-no-kyojin/ymir/pkg/logger"
)

// Logger ...
type Logger struct {
	log logger.Logger
}

// NewLogger ...
func NewLogger(log logger.Logger) *Logger {
	return &Logger{log: log}
}

// Log ...
func (l *Logger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.log.Info(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL))
		next.ServeHTTP(w, r)
	})
}
