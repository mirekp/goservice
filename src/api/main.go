package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/gorilla/mux"
)

const HostName = "localhost"
const Port = 9999
const AccountsRoute = "/accounts"
const MessagesRoute = "/messages"
const MessagesIDRoute = "/messages/{id}"

func main() {
    router := mux.NewRouter()
    router.HandleFunc(AccountsRoute, GetAccountsHandler).Methods("GET")
    router.HandleFunc(AccountsRoute, CreateAccountHandler).Methods("POST")
    router.HandleFunc(MessagesRoute, msgGetAllHandler).Methods("GET")
    router.HandleFunc(MessagesRoute, msgPostHandler).Methods("POST")
    router.HandleFunc(MessagesIDRoute, msgGetHandler).Methods("GET")
    router.HandleFunc(MessagesIDRoute, msgEditHandler).Methods("PATCH")
    router.HandleFunc(MessagesIDRoute, msgDeleteHandler).Methods("DELETE")

    fmt.Printf("Listening on port %d...\n", Port)
    err := http.ListenAndServe(fmt.Sprintf("%s:%d", HostName, Port), router)
    if err != nil {
        log.Fatal("Unable to listen on port: ", err)
    }
}
