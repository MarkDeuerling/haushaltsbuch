package logger

import (
	"fmt"
	"log/slog"
)

// Logger Interface for Logging system wide
type Logger interface {
	Error(string)
	Warning(string)
	Info(string)
	Debug(string)
	Fatal(string)
}

// SystemLogger ...
type SystemLogger struct {
	log *slog.Logger
}

// New SystemLogger
func New() *SystemLogger {
	return &SystemLogger{log: slog.New(slog.Default().Handler())}
}

// Info log
func (s *SystemLogger) Info(msg string) {
	s.log.Info(msg)
}

// Warning log
func (s *SystemLogger) Warning(msg string) {
	s.log.Warn(msg)
}

// Error log
func (s *SystemLogger) Error(msg string) {
	s.log.Error(msg)
}

// Fatal log
func (s *SystemLogger) Fatal(msg string) {
	fmt.Println("Not implemented")
}

// Debug log
func (s *SystemLogger) Debug(msg string) {
	s.log.Debug(msg)
}
