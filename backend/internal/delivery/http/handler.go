package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	stdhttp "net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Nok1o/lab1_rsoi/backend/internal/domain"
	"github.com/Nok1o/lab1_rsoi/backend/internal/usecase/person"
)

const maxRequestBodySize = 1 << 20

type PersonUseCase interface {
	Create(ctx context.Context, input person.CreateInput) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Person, error)
	List(ctx context.Context) ([]domain.Person, error)
	Update(ctx context.Context, id int64, input person.UpdateInput) (domain.Person, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	persons PersonUseCase
	logger  *slog.Logger
}

func NewHandler(persons PersonUseCase, logger *slog.Logger) *Handler {
	return &Handler{persons: persons, logger: logger}
}

func (h *Handler) Create(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var request createPersonRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return
	}

	id, err := h.persons.Create(r.Context(), request.toInput())
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.Header().Set("Location", "/api/v1/persons/"+strconv.FormatInt(id, 10))
	w.WriteHeader(stdhttp.StatusCreated)
}

func (h *Handler) GetByID(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	result, err := h.persons.GetByID(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, toPersonResponse(result))
}

func (h *Handler) List(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	result, err := h.persons.List(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, toPersonResponses(result))
}

func (h *Handler) Update(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var request updatePersonRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, stdhttp.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return
	}

	result, err := h.persons.Update(r.Context(), id, request.toInput())
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, toPersonResponse(result))
}

func (h *Handler) Delete(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.persons.Delete(r.Context(), id); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(stdhttp.StatusNoContent)
}

func (h *Handler) writeError(w stdhttp.ResponseWriter, err error) {
	var validationErr *person.ValidationError
	switch {
	case errors.As(err, &validationErr):
		writeJSON(w, stdhttp.StatusBadRequest, validationErrorResponse{
			Message: validationErr.Error(),
			Errors:  validationErr.Fields,
		})
	case errors.Is(err, domain.ErrPersonNotFound):
		writeJSON(w, stdhttp.StatusNotFound, errorResponse{Message: domain.ErrPersonNotFound.Error()})
	default:
		h.logger.Error("request failed", "error", err)
		writeJSON(w, stdhttp.StatusInternalServerError, errorResponse{Message: "internal server error"})
	}
}

func parseID(w stdhttp.ResponseWriter, r *stdhttp.Request) (int64, bool) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	if err != nil || id <= 0 {
		writeJSON(w, stdhttp.StatusBadRequest, errorResponse{Message: "id must be a positive integer"})
		return 0, false
	}
	return id, true
}

func decodeJSON(w stdhttp.ResponseWriter, r *stdhttp.Request, target any) error {
	r.Body = stdhttp.MaxBytesReader(w, r.Body, maxRequestBodySize)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeJSON(w stdhttp.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
