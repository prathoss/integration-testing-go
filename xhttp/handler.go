package xhttp

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/prathoss/integration_testing/domain"
	"github.com/prathoss/integration_testing/logging"
)

var _ http.Handler = (Handler)(nil)

type Handler func(http.ResponseWriter, *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		handleError(r.Context(), w, err)
	}
}

var _ http.Handler = (HandlerBody[any])(nil)

type HandlerBody[T any] func(w http.ResponseWriter, r *http.Request, body T) error

func (h HandlerBody[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var model T
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h(w, r, model); err != nil {
		handleError(r.Context(), w, err)
	}
}

type ProblemDetailError struct {
	Detail    string `json:"detail"`
	Parameter string `json:"parameter"`
}

type ProblemDetail struct {
	Type   string               `json:"type"`
	Status int                  `json:"status"`
	Title  string               `json:"title"`
	Detail string               `json:"detail"`
	Code   string               `json:"code"`
	Errors []ProblemDetailError `json:"errors"`
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	var notFound *domain.ErrNotFound
	if errors.As(err, &notFound) {
		pd := ProblemDetail{
			Status: http.StatusNotFound,
			Title:  "Not Found",
			Detail: err.Error(),
			Code:   http.StatusText(http.StatusNotFound),
		}
		if err := json.NewEncoder(w).Encode(pd); err != nil {
			slog.Error("could not encode error response to json", logging.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	slog.ErrorContext(ctx, "internal error", logging.Err(err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
