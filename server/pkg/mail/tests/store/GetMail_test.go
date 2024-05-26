package store_test

import (
	"fmt"
	"passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMail(t *testing.T) {
	var (
		store    mail.MailStore
		testMail model.MailEntity
	)

	beforeEach := func() {
		store = mail.NewStore(testDatabase.DB)
	}

	afterEach := func() {
		err = testDatabase.DeleteItemsFromTable("mail")
		fmt.Println("delete table items error", err)
	}

	t.Run("should return nil and error record not found when mail is not found", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		testDatabase.DeleteItemsFromTable("mail")
		badUUID := uuid.New()
		badUser := "badUser"

		// Act
		mail, err := store.GetMail(badUUID, badUser)

		// Assert
		assert.Nil(t, mail)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

	t.Run("should return nil and error record not found when mail is in database but user is not the recipient or sender", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		setupMail := mockMail(1)
		insertMails(setupMail, testDatabase.DB)
		allMails := retrieveMails(testDatabase.DB)
		testMail = allMails[0]
		badUser := "badUser"

		// Act
		mail, err := store.GetMail(testMail.ID, badUser)

		// Assert
		assert.Nil(t, mail)
		assert.Equal(t, "sql: no rows in result set", err.Error())
	})

	t.Run("should return mail when mail is in database and user is the recipient", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		setupMail := mockMail(1)
		insertMails(setupMail, testDatabase.DB)
		allMails := retrieveMails(testDatabase.DB)
		testMail = allMails[0]

		// Act
		mail, err := store.GetMail(testMail.ID, testMail.Recipient)

		// Assert
		assert.Equal(t, testMail, *mail)
		assert.Nil(t, err)
	})

	t.Run("should return mail when mail is in database and user is the sender", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		setupMail := mockMail(1)
		insertMails(setupMail, testDatabase.DB)
		allMails := retrieveMails(testDatabase.DB)
		testMail = allMails[0]

		// Act
		mail, err := store.GetMail(testMail.ID, testMail.Sender)

		// Assert
		assert.Equal(t, testMail, *mail)
		assert.Nil(t, err)
	})

	t.Run("should return error when error is occurred", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		testDatabase.DropTestTable()
		defer testDatabase.CreateTestTable()
		testUser := "testUser"

		// Act
		mail, err := store.GetMail(testMail.ID, testUser)

		// Assert
		assert.Nil(t, mail)
		assert.NotNil(t, err)
	})
}
