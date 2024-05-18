package auth_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"passwordless-mail-server/pkg/auth"
	"passwordless-mail-server/pkg/model"
	"passwordless-mail-server/pkg/util"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testDatabase util.TestDatabase
var err error

func TestMain(m *testing.M) {
	testDatabase, err = util.NewTestDatabase()
	err = testDatabase.CreateTestTable()
	fmt.Println("create table error", err)
	defer testDatabase.DropTestTable()
	code := m.Run()
	os.Exit(code)
}

func TestStore_GetUsedUUID(t *testing.T) {
	var store auth.UuidStore

	beforeEach := func() {
		store = auth.NewUUIDStore(testDatabase.DB)
	}

	afterEach := func() {
		err = testDatabase.DeleteItemsFromTable("used_uuid")
		fmt.Println("drop table error", err)
	}

	t.Run("should return used uuid when uuid is stored", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockUUIDs := mockUsedUUID(1)
		insertUsedUUIDs(mockUUIDs, testDatabase.DB)
		uuid := mockUUIDs[0].UUID

		// Act
		entity, err := store.GetUsedUUID(uuid)

		// Assert
		assert.Equal(t, nil, err)
		assert.Equal(t, uuid, entity.UUID)
	})

	t.Run("should return no uuid value and no error when uuid is not stored", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		notStoredUUID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
		uuid, parseErr := uuid.Parse(notStoredUUID)

		// Act
		entity, queryErr := store.GetUsedUUID(uuid)

		// Assert
		assert.Equal(t, nil, parseErr)
		assert.Equal(t, nil, queryErr)
		assert.Nil(t, entity)
	})

	t.Run("should return error when error occurred", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockUUIDs := mockUsedUUID(1)
		insertUsedUUIDs(mockUUIDs, testDatabase.DB)
		uuid := mockUUIDs[0].UUID

		// Act
		testDatabase.DropTestTable() // drop table to force error
		defer testDatabase.CreateTestTable()
		entity, err := store.GetUsedUUID(uuid)

		// Assert
		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestStore_InsertUsedUUID(t *testing.T) {
	var store auth.UuidStore

	beforeEach := func() {
		testDatabase, err = util.NewTestDatabase()
		store = auth.NewUUIDStore(testDatabase.DB)
	}

	afterEach := func() {
		err = testDatabase.DropTestTable()
		fmt.Println("drop table error", err)
		err = testDatabase.CreateTestTable()
		fmt.Println("create table error", err)
	}

	t.Run("should insert used uuid to database", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		uuid := uuid.New()

		// Act
		insertErr := store.InsertUsedUUID(uuid)

		// Assert
		usedUUIDs := retrieveUsedUUIDs(testDatabase.DB)
		fmt.Println("usedUUIDs", usedUUIDs)
		assert.Equal(t, nil, insertErr)
		assert.Equal(t, 1, len(usedUUIDs))
		assert.Equal(t, uuid, usedUUIDs[0].UUID)
	})

	t.Run("should return error when uuid is already stored", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		uuid := uuid.New()
		insertUsedUUIDs([]model.UsedUUIDEntity{{UUID: uuid}}, testDatabase.DB)

		// Act
		insertErr := store.InsertUsedUUID(uuid)

		// Assert
		usedUUIDs := retrieveUsedUUIDs(testDatabase.DB)
		expectErr := fmt.Errorf("uuid %s is already used", uuid)
		assert.Equal(t, expectErr, insertErr)
		assert.Equal(t, 1, len(usedUUIDs))
		assert.Equal(t, uuid, usedUUIDs[0].UUID)
	})

	t.Run("should return error when error occurred", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		uuid := uuid.New()

		// Act
		testDatabase.DropTestTable() // drop table to force error
		defer testDatabase.CreateTestTable()
		insertErr := store.InsertUsedUUID(uuid)

		// Assert
		usedUUIDs := retrieveUsedUUIDs(testDatabase.DB)
		assert.Error(t, insertErr)
		assert.Equal(t, 0, len(usedUUIDs))
	})
}

func mockUsedUUID(amount int) []model.UsedUUIDEntity {
	if amount < 1 || amount > 100 {
		log.Fatal("amount should be between 1 and 99, heehee! ow!")
	}
	uuids := []model.UsedUUIDEntity{}
	fixPrefix := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa"
	for i := 1; i < amount+1; i++ {
		// if 1 digit, add 0 in front
		suffix := fmt.Sprintf("%02d", i)
		genUUID, err := uuid.Parse(fixPrefix + suffix)
		if err != nil {
			log.Fatal(err)
		}
		uuids = append(uuids, model.UsedUUIDEntity{
			UUID: genUUID,
		})
	}
	return uuids
}

func insertUsedUUIDs(uuids []model.UsedUUIDEntity, db *sql.DB) error {
	queryScript := "INSERT INTO used_uuid (uuid) VALUES ($1) ON CONFLICT (uuid) DO NOTHING"
	for _, uuid := range uuids {
		_, err := db.Exec(queryScript, uuid.UUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveUsedUUIDs(db *sql.DB) []model.UsedUUIDEntity {
	queryScript := "SELECT * FROM used_uuid"
	var uuids []model.UsedUUIDEntity
	rows, err := db.Query(queryScript)
	if err != nil {
		return []model.UsedUUIDEntity{}
	}
	defer rows.Close()
	for rows.Next() {
		var uuid model.UsedUUIDEntity
		err := rows.Scan(&uuid.UUID, &uuid.CreatedAt)
		if err != nil {
			return []model.UsedUUIDEntity{}
		}
		uuids = append(uuids, uuid)
	}
	return uuids
}
