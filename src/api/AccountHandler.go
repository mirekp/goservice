package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
)

func GetAccountsHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := getAccountIDFromRequest(request)

	if err != nil {
		fmt.Printf("Unable to extract ID from request\n")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := GetAccountByID(id)

	if err != nil {
		fmt.Printf("Unable to find requested account\n")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(writer).Encode(serialiseAccountToJSON(account))
}

func CreateAccountHandler(writer http.ResponseWriter, request *http.Request) {

	// parse the JSON request
	var accountJSON AccountJSON
	err := json.NewDecoder(request.Body).Decode(&accountJSON)
	accountData := accountJSON.AccountData

	if err != nil || accountData.Type != "account" || accountData.Attr.Username == "" || accountData.Attr.Password == "" {
		fmt.Printf("Unable to process your request\n")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = GetAccountByName(accountData.Attr.Username)

	if err == nil {
		fmt.Printf("Account already exits\n")
		writer.WriteHeader(http.StatusConflict)
		return
	}

	account, err := CreateAccount(accountData.Attr.Username, accountData.Attr.Password)

	if err != nil {
		fmt.Printf("Unable to create account due to internal error\n")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Location", fmt.Sprintf("http://%s:%d%s/%d", HostName, Port, RouteAccounts, account.id))
	writer.WriteHeader(http.StatusCreated)
	createdAccount := serialiseAccountToJSON(account)
	json.NewEncoder(writer).Encode(createdAccount)
	return
}

func getAccountIDFromRequest(request *http.Request) (int, error) {
	vars := mux.Vars(request)
	stringID, ok := vars["id"]

	if ok == false {
		return 0, fmt.Errorf("No message-id provided")
	}

	id, err := strconv.Atoi(stringID)
	if err != nil {
		return 0, fmt.Errorf("Bad msg-id format")
	}

	return id, nil
}

