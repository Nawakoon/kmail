package auth

import (
	"database/sql"
	"errors"
	"passwordless-mail-server/pkg/model"

	"github.com/google/uuid"
)

type UuidStore interface {
	GetUsedUUID(uuid uuid.UUID) (*model.UsedUUIDEntity, error)
	InsertUsedUUID(uuid uuid.UUID) error
}

type Store struct {
	db *sql.DB
}

func NewUUIDStore(database *sql.DB) UuidStore {
	return &Store{
		db: database,
	}
}

// not found 		-> nil, nil (uuid is not used)
// found 			-> entity, nil (uuid is used)
// side effect err	-> nil, error (error occurred)
func (s *Store) GetUsedUUID(uuid uuid.UUID) (*model.UsedUUIDEntity, error) {
	queryScript := "SELECT * FROM used_uuid WHERE uuid = $1"
	row := s.db.QueryRow(queryScript, uuid)
	entity := &model.UsedUUIDEntity{}
	err := row.Scan(&entity.UUID, &entity.CreatedAt)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// success 			-> no error
// duplicate uuid 	-> error 'uuid ... is already used'
// side effect err	-> error <error details>
func (s *Store) InsertUsedUUID(uuid uuid.UUID) error {
	queryScript := "INSERT INTO used_uuid (uuid) VALUES ($1)"
	_, err := s.db.Exec(queryScript, uuid)
	if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"used_uuid_pkey\"" {
		return errors.New("uuid " + uuid.String() + " is already used")
	}
	if err != nil {
		return err
	}

	return nil
}
