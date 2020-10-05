package Controller

import (
	"Edgex-Export_Receiver/EdgexData"
	mux "github.com/gorilla/mux"
	"net/http"
)

func InitRestRoutes() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/push", EdgexData.Receive).Methods(http.MethodPost)
	s.HandleFunc("/pull", EdgexData.Pull).Methods(http.MethodGet)
	return r
}