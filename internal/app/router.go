package app

import (
	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmgorilla/v2"
	"net/http"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()

	apmgorilla.Instrument(r)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Changa!"))
	}).Methods(http.MethodGet)

	return r
}
