package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mandoju/BindAPI/handlers/login"
	"github.com/mandoju/BindAPI/utils/Database"
	"github.com/mandoju/BindAPI/handlers/domain"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "get called"}`))
}

func main() {
	Database.InitializeDb()
	defer Database.Db.Close()
	r := mux.NewRouter()
	//api := r.PathPrefix("/").Subrouter()
	r.HandleFunc("/login", login.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/refresh", login.RefreshHandler).Methods(http.MethodPost)
	r.HandleFunc("/register", login.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/domains", domain.GetDomainHandler).Methods(http.MethodGet)
	r.HandleFunc("/", get).Methods(http.MethodGet)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err.Error())
	}
}
