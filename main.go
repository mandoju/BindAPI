package main

import (
	"github.com/mandoju/BindAPI/handlers"
	"net/http"
	"github.com/gorilla/mux"
)

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func main() {
	r := mux.NewRouter()
	//api := r.PathPrefix("/").Subrouter()
	r.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/", get).Methods(http.MethodGet)
	http.ListenAndServe(":8080", r)

}
