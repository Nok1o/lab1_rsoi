package http

import (
	stdhttp "net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) stdhttp.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/persons", handler.List).Methods(stdhttp.MethodGet)
	api.HandleFunc("/persons", handler.Create).Methods(stdhttp.MethodPost)
	api.HandleFunc("/persons/{id}", handler.GetByID).Methods(stdhttp.MethodGet)
	api.HandleFunc("/persons/{id}", handler.Update).Methods(stdhttp.MethodPatch)
	api.HandleFunc("/persons/{id}", handler.Delete).Methods(stdhttp.MethodDelete)

	router.HandleFunc("/health", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods(stdhttp.MethodGet)

	return router
}
