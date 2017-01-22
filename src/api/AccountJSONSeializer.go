package main

import "fmt"

type AccountJSON struct {
	AccountData AccountData  `json:"data"`
}

type AccountData struct {
	ID   string       `json:"id,omitempty"`
	Type string       `json:"type"`
	Attr AccountAttr  `json:"attributes"`
}

type AccountAttr struct {
	Username string       `json:"username"`
	Password string       `json:"password,omitempty"`
}


func serialiseAccountToJSON(account Account) AccountJSON {
	return AccountJSON{
		AccountData: AccountData{
			ID: fmt.Sprintf("%d", account.id),
			Type: "account",
			Attr: AccountAttr{
				Username: account.Username,
			},
		},
	}
}
