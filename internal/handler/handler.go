package handler

import (
	"github.com/gorilla/mux"
	"log-aggregator/internal/service"
	"net/http"
)

type Handler struct {
	Service service.IService
}

func NewHandler(service service.IService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) SetupRoutes(r *mux.Router) {

	public := r.NewRoute().Subrouter()

	// user
	public.HandleFunc("/event/save", h.SaveEvent).Methods(http.MethodPost)
}
