package main

import (
	"github.com/gorilla/mux"
	"github.com/mandoju/BindAPI/handlers/login"
	"net/http"
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
	r.HandleFunc("/login", login.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/", get).Methods(http.MethodGet)
	http.ListenAndServe(":8080", r)

}
