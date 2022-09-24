package main

import (
	"chat-server/app"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user", app.CreateUser).Methods("POST")
	r.HandleFunc("/user/login", app.LoginUser).Methods("POST")

	log.Println("Listening ...")
	http.ListenAndServe(":8080", r)
}
