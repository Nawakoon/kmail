package model

import "github.com/google/uuid"

type Mail struct {
	ID      uuid.UUID `json:"id"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
	Body    string    `json:"body"`
}

type MailFileContent struct {
	To      *string `json:"to"`
	Subject *string `json:"subject"`
	Body    *string `json:"body`
}
