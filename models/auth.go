package models

// Authentication used to parse Authentication requests
type Authentication struct {
	Login    string `schema:"login"`
	Password string `schema:"password"`
}
