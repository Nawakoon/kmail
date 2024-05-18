package model

import "github.com/google/uuid"

type Mail struct {
	ID      uuid.UUID `json:"id"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
	Body    string    `json:"body"`
}

type InboxResponse struct {
	Inbox []Mail `json:"inbox"`
	Total int    `json:"total"`
}

type SendMailResponse struct {
	ID uuid.UUID `json:"id"`
}

// SQL table schema

type MailEntity struct {
	ID          uuid.UUID `db:"id"`
	Recipient   string    `db:"recipient"`
	Sender      string    `db:"sender"`
	MailSubject string    `db:"mail_subject"`
	Body        string    `db:"body"`
	SentAt      string    `db:"sent_at"`
}

type UsedUUIDEntity struct {
	UUID      uuid.UUID `db:"uuid"`
	CreatedAt string    `db:"created_at"`
}
