package store_test

import (
	"fmt"
	"passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertMail(t *testing.T) {
	var store mail.MailStore

	beforeEach := func() {
		store = mail.NewStore(testDatabase.DB)
	}

	afterEach := func() {
		err = testDatabase.DeleteItemsFromTable("mail")
		fmt.Println("delete table items error", err)
	}

	t.Run("should insert mail into the database", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		testMail := model.Mail{
			From:    "sender public key",
			To:      "recipient public key",
			Subject: "test subject",
			Body:    "test body",
		}
		preStoredMails := retrieveMails(testDatabase.DB)

		// Act
		insertedMail, err := store.InsertMail(testMail)
		fmt.Println("insert mail error", err)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, insertedMail.ID)
		postStoredMails := retrieveMails(testDatabase.DB)
		assert.Equal(t, len(preStoredMails)+1, len(postStoredMails))
		assert.Equal(t, testMail.From, postStoredMails[0].Sender)
		assert.Equal(t, testMail.To, postStoredMails[0].Recipient)
		assert.Equal(t, testMail.Subject, postStoredMails[0].MailSubject)
		assert.Equal(t, testMail.Body, postStoredMails[0].Body)
	})

	t.Run("should return error when error is occurred", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		testMail := model.Mail{
			From:    "sender public key",
			To:      "recipient public key",
			Subject: "test subject",
			Body:    "test body",
		}
		preStoredMails := retrieveMails(testDatabase.DB)

		// Act
		testDatabase.DropTestTable()
		defer testDatabase.CreateTestTable()
		insertedMail, err := store.InsertMail(testMail)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, insertedMail)
		postStoredMails := retrieveMails(testDatabase.DB)
		assert.Equal(t, len(preStoredMails), len(postStoredMails))
	})
}
