package store_test

import (
	"fmt"
	"passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalMailsReceived(t *testing.T) {
	var store mail.MailStore

	beforeEach := func() {
		store = mail.NewStore(testDatabase.DB)
	}

	afterEach := func() {
		err = testDatabase.DeleteItemsFromTable("mail")
		fmt.Println("delete table items error", err)
	}

	t.Run("should return 0 when no mail is in the inbox", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockMails := mockMail(2)
		insertMails(mockMails, testDatabase.DB)

		// Act
		totalMailsReceived, err := store.GetTotalMailsReceived("random-recipient")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, totalMailsReceived)
	})

	t.Run("should return total mails received when mails are in the inbox", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		var mails []model.MailEntity
		mailAmount := 15

		for i := 0; i < mailAmount; i++ {
			mail := model.MailEntity{
				Recipient:   "recipient-1",
				Sender:      fmt.Sprintf("sender-%d", i+1),
				MailSubject: fmt.Sprintf("subject-%d", i+1),
				Body:        fmt.Sprintf("body-%d", i+1),
			}
			mails = append(mails, mail)
		}
		insertMails(mails, testDatabase.DB)

		// Act
		totalMailsReceived, err := store.GetTotalMailsReceived("recipient-1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, mailAmount, totalMailsReceived)
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		// Arrange
		beforeEach()
		defer afterEach()
		mockMails := mockMail(2)
		insertMails(mockMails, testDatabase.DB)

		// Act
		testDatabase.DropTestTable()
		defer testDatabase.CreateTestTable()
		totalMailsReceived, err := store.GetTotalMailsReceived("random-recipient")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, 0, totalMailsReceived)
	})
}
