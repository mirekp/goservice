package model

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	HashedPassword  string   `json:"hashedpassword,omitempty"`
}

var People []Person
