package mail

import (
	"database/sql"
	"passwordless-mail-server/pkg/model"

	"github.com/google/uuid"
)

type StoreGetInboxQuery struct {
	Recipient string
	Offset    int
	Limit     int
}

type MailStore interface {
	GetInbox(query StoreGetInboxQuery) ([]model.MailEntity, error)
	GetMail()
	InsertMail(mail model.Mail) (*model.MailEntity, error)
}

type Store struct {
	db *sql.DB
}

func NewStore(database *sql.DB) MailStore {
	return &Store{
		db: database,
	}
}

func (s *Store) GetInbox(query StoreGetInboxQuery) ([]model.MailEntity, error) {
	getInboxQuery := `
		SELECT * FROM mail
		WHERE recipient = $1
		LIMIT $2
		OFFSET $3
	`

	rows, err := s.db.Query(getInboxQuery, query.Recipient, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var inbox []model.MailEntity
	for rows.Next() {
		var mail model.MailEntity
		err := rows.Scan(&mail.ID, &mail.Recipient, &mail.Sender, &mail.MailSubject, &mail.Body, &mail.SentAt)
		if err != nil {
			return nil, err
		}
		inbox = append(inbox, mail)
	}

	return inbox, nil
}

func (s *Store) GetMail() {
}

func (s *Store) InsertMail(mail model.Mail) (*model.MailEntity, error) {
	queryScript := `
		INSERT INTO mail (recipient, sender, mail_subject, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var mailId uuid.UUID
	err := s.db.QueryRow(
		queryScript,
		mail.To,
		mail.From,
		mail.Subject,
		mail.Body,
	).Scan(&mailId)

	if err != nil {
		return nil, err
	}

	return &model.MailEntity{
		ID:          mailId,
		Recipient:   mail.To,
		Sender:      mail.From,
		MailSubject: mail.Subject,
		Body:        mail.Body,
	}, nil
}
