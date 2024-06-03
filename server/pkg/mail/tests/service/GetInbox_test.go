package service_test

import (
	"passwordless-mail-client/pkg/account"
	"passwordless-mail-client/pkg/request"
	authmocks "passwordless-mail-server/pkg/auth/mocks"
	"passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/mail/mocks"
	mailmock "passwordless-mail-server/pkg/mail/mocks"
	"passwordless-mail-server/pkg/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInbox(t *testing.T) {

	const TestPrivateKey = "35c03d4a383c899345cc2e8d49417a92b7654fab37d404783dac84e3fcf5d66e"

	var (
		testAccount *account.Account
		err         error
		testUUID    uuid.UUID

		mockMailStore        mailmock.MailStore
		mockUUIDStore        authmocks.UuidStore
		mailService          mail.MailService
		resMailStoreGetInbox []model.MailEntity
		errMailStoreGetInbox error
		resMailStoreGetTotal int
		errMailStoreGetTotal error
		resUUIDStoreGetUUID  *model.UsedUUIDEntity
		errUUIDStoreGetUUID  error
		errInsertUsedUUID    error
	)

	beforeEach := func() {
		// setup test account
		testAccount, err = account.ConnectAccount(TestPrivateKey)
		if err != nil || testAccount == nil {
			t.Errorf("failed to connect account")
		}

		testUUID, err = uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
		assert.NoError(t, err)

		// setup mail service
		mockMailStore = mocks.MailStore{}
		mockUUIDStore = authmocks.UuidStore{}
		mailService = mail.NewService(&mockMailStore, &mockUUIDStore)

		resMailStoreGetInbox = []model.MailEntity{
			{
				ID:          testUUID,
				Recipient:   "recipient",
				Sender:      "sender",
				MailSubject: "mail subject",
				Body:        "mail body",
				SentAt:      "2021-01-01 00:00:00",
			},
		}
		errMailStoreGetInbox = nil
		resUUIDStoreGetUUID = nil
		errUUIDStoreGetUUID = nil
		errInsertUsedUUID = nil

		mockMailStore.On("GetInbox", mock.Anything).Return(resMailStoreGetInbox, errMailStoreGetInbox)
		mockMailStore.On("GetTotalMailsReceived", mock.Anything).Return(resMailStoreGetTotal, errMailStoreGetTotal)
		mockUUIDStore.On("GetUsedUUID", mock.Anything).Return(resUUIDStoreGetUUID, errUUIDStoreGetUUID)
		mockUUIDStore.On("InsertUsedUUID", mock.Anything).Return(errInsertUsedUUID)

	}

	t.Run("should have correct inbox query", func(t *testing.T) {
		// Arrange
		beforeEach()
		message, newMsgErr := request.NewGetInbox()
		signedMassage, signErr := testAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMassage,
		}
		publicKey := testAccount.PublicKey
		serviceGetInboxQuery := mail.ServiceGetInboxQuery{
			Recipient: testAccount.GetAddress(),
			Page:      3,
			Limit:     10,
		}

		// Act
		mailService.GetInbox(requestBody, publicKey, serviceGetInboxQuery)

		// Assert
		assert.NoError(t, newMsgErr)
		assert.NoError(t, signErr)
		expectedStoreQuery := mail.StoreGetInboxQuery{
			Recipient: testAccount.GetAddress(),
			Offset:    20,
			Limit:     10,
		}
		mockMailStore.AssertCalled(t, "GetInbox", expectedStoreQuery)
	})
}
