package main

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type Account struct {
	id            	int
	Username      	string
	HashedPassword  string
}

// data storage
var accounts = []Account{}
var lastUserID = 0

func CreateAccount(username, password string) (Account, error) {

	// check if a given username is available
	_, err := GetAccountByName(username)
	if err == nil {
		return Account{}, fmt.Errorf("Account already exists")
	}

	// calculate hash of the password as we don't want to store plaintext passwords in our model
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}

	fmt.Printf("Creating new account %s\n", username)

	account := Account{
		id: lastUserID,
		Username: username,
		HashedPassword: string(hashedPassword),
	}

	accounts = append(accounts, account)
	lastUserID++

	return account, nil
}

func GetAccountByName(username string) (Account, error) {

	for _, account := range accounts {
		if account.Username == username {
			return account, nil
		}
	}

	return Account{}, fmt.Errorf("Unable to find the account")
}

func GetAccountByID(id int) (Account, error) {

	for _, account := range accounts {
		if account.id == id {
			return account, nil
		}
	}

	return Account{}, fmt.Errorf("Unable to find the account")
}
