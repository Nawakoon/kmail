package mail

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"passwordless-mail-client/pkg/account"
	"passwordless-mail-client/pkg/request"
	"passwordless-mail-server/pkg/auth"
	"passwordless-mail-server/pkg/model"
	"time"
)

type ServiceGetInboxQuery struct {
	Recipient string
	Page      int
	Limit     int
}

type MailService interface {
	GetInbox(request model.RequestBody, publicKey *ecdsa.PublicKey, query ServiceGetInboxQuery) (model.InboxResponse, error)
	GetMail()
	SendMail(request model.RequestBody, publicKey *ecdsa.PublicKey) (model.SendMailResponse, error)
}

type Service struct {
	mailStore MailStore
	uuidStore auth.UuidStore
}

func NewService(mailStore MailStore, uuidStore auth.UuidStore) MailService {
	return &Service{
		mailStore: mailStore,
		uuidStore: uuidStore,
	}
}

func (s *Service) GetInbox(
	requestBody model.RequestBody,
	publicKey *ecdsa.PublicKey,
	query ServiceGetInboxQuery,
) (model.InboxResponse, error) {
	const TIMEOUT = time.Duration(3 * time.Minute)

	isVerify := account.Verify(
		publicKey,
		[]byte(requestBody.Data),
		[]byte(requestBody.Signature),
	)

	if !isVerify {
		return model.InboxResponse{}, fmt.Errorf("validation failed")
	}

	var message request.GetInboxRequest
	err := json.Unmarshal([]byte(requestBody.Data), &message)
	if err != nil {
		return model.InboxResponse{}, fmt.Errorf("bad request")
	}

	// turn message.Timestamp string into time.Time
	timestamp, err := time.Parse(time.RFC3339, message.Timestamp)
	if err != nil {
		return model.InboxResponse{}, fmt.Errorf("bad request")
	}
	isTimeout := time.Since(timestamp) > TIMEOUT
	if isTimeout {
		return model.InboxResponse{}, fmt.Errorf("message timeout")
	}

	usedUUID, err := s.uuidStore.GetUsedUUID(message.ID)
	if err != nil {
		return model.InboxResponse{}, err
	}
	if usedUUID != nil {
		return model.InboxResponse{}, fmt.Errorf("uuid is already used")
	}
	err = s.uuidStore.InsertUsedUUID(message.ID)
	if err != nil {
		return model.InboxResponse{}, err
	}

	storeQuery := StoreGetInboxQuery{
		Recipient: query.Recipient,
		Limit:     query.Limit,
		Offset:    (query.Page - 1) * query.Limit,
	}

	inbox, err := s.mailStore.GetInbox(storeQuery)
	if err != nil {
		return model.InboxResponse{}, err
	}

	parsedInbox := []model.Mail{}
	for _, mailEntity := range inbox {
		parsedInbox = append(parsedInbox, model.Mail{
			ID:      mailEntity.ID,
			From:    mailEntity.Sender,
			To:      mailEntity.Recipient,
			Subject: mailEntity.MailSubject,
			Body:    mailEntity.Body,
		})
	}

	inboxResponse := model.InboxResponse{
		Inbox: parsedInbox,
		Total: len(parsedInbox),
	}
	return inboxResponse, nil
}

func (s *Service) GetMail() {
}

func (s *Service) SendMail(
	requestBody model.RequestBody,
	publicKey *ecdsa.PublicKey,
) (model.SendMailResponse, error) {
	const TIMEOUT = time.Duration(3 * time.Minute)

	isVerify := account.Verify(
		publicKey,
		[]byte(requestBody.Data),
		[]byte(requestBody.Signature),
	)
	if !isVerify {
		return model.SendMailResponse{}, fmt.Errorf("validation failed")
	}

	var message request.SendEmailRequest
	err := json.Unmarshal([]byte(requestBody.Data), &message)
	if err != nil {
		return model.SendMailResponse{}, fmt.Errorf("bad request")
	}

	// check if recipient is a valid public key
	_, err = account.HexToPublicKey(message.Recipient)
	if err != nil {
		return model.SendMailResponse{}, fmt.Errorf("bad request")
	}

	timestamp, err := time.Parse(time.RFC3339, message.Timestamp)
	if err != nil {
		return model.SendMailResponse{}, fmt.Errorf("bad request")
	}
	isTimeout := time.Since(timestamp) > TIMEOUT
	if isTimeout {
		return model.SendMailResponse{}, fmt.Errorf("message timeout")
	}

	usedUUID, err := s.uuidStore.GetUsedUUID(message.ID)
	if err != nil {
		return model.SendMailResponse{}, err
	}
	if usedUUID != nil {
		return model.SendMailResponse{}, fmt.Errorf("uuid is already used")
	}
	err = s.uuidStore.InsertUsedUUID(message.ID)
	if err != nil {
		return model.SendMailResponse{}, err
	}

	sender := account.PublicKeyToHex(publicKey)

	insertedMail, err := s.mailStore.InsertMail(model.Mail{
		From:    sender,
		To:      message.Recipient,
		Subject: message.Subject,
		Body:    message.Body,
	})
	if err != nil {
		return model.SendMailResponse{}, err
	}

	return model.SendMailResponse{
		ID: insertedMail.ID,
	}, nil
}
