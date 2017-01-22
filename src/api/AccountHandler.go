package main

import (
    "net/http"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
    "fmt"
)

/*

Sample account JSON:

{
    data : {
        "id": "username",
        "type": "account",
        "attributes": {
            "password", "plaintext-password"
         }
    }
}
 */

type AccountJSON struct {
    AccountData AccountData `json:"data"`
}

type AccountData struct {
    ID             string   `json:"id"`
    Type           string   `json:"type"`
    Attr           userAttr `json:"attributes"`
}

type userAttr struct {
    Password       string   `json:"password,omitempty"`
}

var accounts = make(map[string]string)

func GetAccountsHandler (writer http.ResponseWriter, _ *http.Request) {

    json.NewEncoder(writer).Encode(accounts)
}

func CreateAccountHandler (writer http.ResponseWriter, request *http.Request) {

    // parse the JSON request
    var accountJSON AccountJSON
    err := json.NewDecoder(request.Body).Decode(&accountJSON)
    accountData := accountJSON.AccountData

    if err != nil || accountData.Type != "account" || accountData.ID == "" || accountData.Attr.Password == "" {
        fmt.Printf("Unable to process your request\n")
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    // stop if a given username already exists in the database
    _, ok := accounts[accountData.ID]
    if ok == true {
        fmt.Printf("User %s already exists\n", accountData.ID)
        writer.WriteHeader(http.StatusConflict)
        return
    }

    // calculate hash of the password as we don't want to store plaintext passwords in our model
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(accountData.Attr.Password), bcrypt.DefaultCost)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        return
    }

    fmt.Println("Creating user:", accountData.ID)
    accounts[accountData.ID] = string(hashedPassword)
}
