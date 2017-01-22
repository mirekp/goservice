package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

const HostName = "localhost"
const Port = 9999

const RouteAccounts = "/accounts"
const RouteAccountsWithID = "/accounts/{id}"
const RouteMessage = "/messages"
const RouteMessageWithID = "/messages/{id}"

func main() {
	router := mux.NewRouter()
	router.HandleFunc(RouteAccountsWithID, GetAccountsHandler).Methods("GET")
	router.HandleFunc(RouteAccounts, CreateAccountHandler).Methods("POST")
	router.HandleFunc(RouteMessage, msgGetAllHandler).Methods("GET")
	router.HandleFunc(RouteMessage, msgPostHandler).Methods("POST")
	router.HandleFunc(RouteMessageWithID, msgGetHandler).Methods("GET")
	router.HandleFunc(RouteMessageWithID, msgEditHandler).Methods("PATCH")
	router.HandleFunc(RouteMessageWithID, msgDeleteHandler).Methods("DELETE")

	fmt.Printf("Listening on %s port %d...\n", HostName, Port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", HostName, Port), router)
	if err != nil {
		log.Fatal("Unable to listen on port: ", err)
	}
}
