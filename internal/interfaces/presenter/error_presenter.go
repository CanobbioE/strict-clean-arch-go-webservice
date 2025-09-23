// Package presenter formats output representation for the external world (JSON payload shape, view models).
// Keeps controllers thin.
package presenter

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
)

// ErrorPresenter prepares an error to be returned to an http interface.
type ErrorPresenter struct {
	logger *slog.Logger
}

// NewErrorPresenter creates a new instance of ErrorPresenter.
func NewErrorPresenter(logger *slog.Logger) *ErrorPresenter {
	return &ErrorPresenter{logger: logger}
}

// Present writes the JSON representation of the error directly to w.
// Additionally, it writes the correct error code to w.
// Internal errors are caught and replaced with a default message.
// If the error is one of the known types, the code is overwritten with the correct one.
func (p *ErrorPresenter) Present(w http.ResponseWriter, err error, code int) {
	switch {
	case err == nil:
		return
	case errors.Is(err, domain.ErrBookNotFound):
		code = http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidBookID):
		code = http.StatusBadRequest
	case code == http.StatusInternalServerError:
		p.handleInternalError(w, err)
		return
	case http.StatusText(code) != "":
	default:
		p.handleInternalError(w, err)
		return
	}
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusText(code),
		"message": err.Error(),
	})
	if err != nil {
		p.logger.With("error", err).Error("failed to write error response")
	}
}

func (p *ErrorPresenter) handleInternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	err = json.NewEncoder(w).Encode(map[string]any{
		"status":  http.StatusText(http.StatusInternalServerError),
		"message": "internal server error",
	})
	if err != nil {
		p.logger.With("error", err).Error("failed to write error response")
	}
}
