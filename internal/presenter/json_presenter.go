package presenter

import (
	"encoding/json"
	"net/http"
)

type presenter interface {
	Successful(data any) error
	Failed(string) error
}

// JSONPresenter is the presenter for the JSON response.
type JSONPresenter struct {
	w http.ResponseWriter
}

// NewJSONPresenter creates a new JSON presenter.
func NewJSONPresenter(w http.ResponseWriter) *JSONPresenter {
	return &JSONPresenter{w: w}
}

// Successful sends a successful JSON response.
func (p *JSONPresenter) Successful(data any) error {
	return json.NewEncoder(p.w).Encode(map[string]any{"data": data})
}

// Failed sends a failed JSON response.
func (p *JSONPresenter) Failed(msg string) error {
	return json.NewEncoder(p.w).Encode(map[string]string{"error": msg})
}
