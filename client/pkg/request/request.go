package request

import (
	"encoding/json"
	"passwordless-mail-client/pkg/model"
	"time"

	"github.com/google/uuid"
)

type ActionName string

const (
	GetInbox  ActionName = "get inbox"
	GetEmail  ActionName = "get email"
	SendEmail ActionName = "send email"
)

type GetInboxRequest struct {
	ID        uuid.UUID `json:"id"`
	Timestamp string    `json:"timestamp"`
}

type GetInboxResponse struct {
	Inbox []model.Mail `json:"inbox"`
	Total int          `json:"total"`
}

type GetEmailRequest struct {
	ID        uuid.UUID `json:"id"`
	Timestamp string    `json:"timestamp"`
	EmailID   uuid.UUID `json:"email_id"`
}

type SendEmailRequest struct {
	ID        uuid.UUID `json:"id"`
	Timestamp string    `json:"timestamp"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject`
	Body      string    `json:"body"`
}

type SendMailResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewGetInbox() ([]byte, error) {
	getInbox := GetInboxRequest{
		ID:        uuid.New(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	inbox, err := json.Marshal(getInbox)
	if err != nil {
		return nil, err
	}

	return inbox, nil
}

func NewGetEmail() ([]byte, error) {
	getEmail := GetEmailRequest{
		ID:        uuid.New(),
		Timestamp: time.Now().Format(time.RFC3339),
		EmailID:   uuid.New(),
	}

	strJSON, err := json.Marshal(getEmail)
	if err != nil {
		return nil, err
	}

	return strJSON, nil
}

func NewSendEmail(recipient string, subject string, body string) ([]byte, error) {
	mail := SendEmailRequest{
		ID:        uuid.New(),
		Timestamp: time.Now().Format(time.RFC3339),
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
	}

	strJSON, err := json.Marshal(mail)
	if err != nil {
		return nil, err
	}

	return strJSON, nil
}
