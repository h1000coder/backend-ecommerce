package model

type Payment struct {
	ID int    `json:"id"`
	Email string `json:"email"`
	Telephone string `json:"telephone"`
	Address string `json:"address"`
}