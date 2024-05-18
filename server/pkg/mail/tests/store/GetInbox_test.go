package store_test

import (
	"fmt"
	"passwordless-mail-server/pkg/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInbox(t *testing.T) {
	var store mail.MailStore

	beforeEach := func() {
		store = mail.NewStore(testDatabase.DB)
	}

	afterEach := func() {
		err = testDatabase.DeleteItemsFromTable("mail")
		fmt.Println("delete table items error", err)
	}

	t.Run("should return empty array of inbox when no mail is in the inbox", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockMails := mockMail(2)
		insertMails(mockMails, testDatabase.DB)

		// Act
		query := mail.StoreGetInboxQuery{
			Recipient: "random-recipient",
			Offset:    0,
			Limit:     10,
		}
		inbox, err := store.GetInbox(query)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(inbox))
	})

	t.Run("should return array of inbox when mails are in the inbox", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockMails := mockMail(3)
		mockMails[0].Recipient = "recipient-1"
		mockMails[1].Recipient = "recipient-1"
		mockMails[2].Recipient = "recipient-2"
		insertMails(mockMails, testDatabase.DB)

		// Act
		query := mail.StoreGetInboxQuery{
			Recipient: "recipient-1",
			Offset:    0,
			Limit:     10,
		}
		inbox, err := store.GetInbox(query)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(inbox))
	})

	t.Run("should return inbox with correct limit and page", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockMails := mockMail(10)
		mockMails[0].Recipient = "recipient-1"
		mockMails[1].Recipient = "recipient-1"
		mockMails[2].Recipient = "recipient-1"
		insertMails(mockMails, testDatabase.DB)

		// Act
		query1 := mail.StoreGetInboxQuery{
			Recipient: "recipient-1",
			Offset:    0,
			Limit:     5,
		}
		query2 := mail.StoreGetInboxQuery{
			Recipient: "recipient-1",
			Offset:    5,
			Limit:     5,
		}
		inbox1, err1 := store.GetInbox(query1)
		inbox2, err2 := store.GetInbox(query2)

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 3, len(inbox1))
		assert.Equal(t, 0, len(inbox2))
	})
}
