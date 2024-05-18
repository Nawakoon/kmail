package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"passwordless-mail-client/pkg/account"
	"passwordless-mail-client/pkg/request"
	"passwordless-mail-server/pkg/model"
	"passwordless-mail-server/pkg/util"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SendMailTestCases(t *testing.T) {

	const (
		SendMailPath    = "http://localhost:8080/mail/send"
		InboxPath       = "http://localhost:8080/mail/inbox"
		TestPrivateKey1 = "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247"
		TestPrivateKey2 = "fd778940ddae63e19e5d2a05604a4d0eaec18b977801299a7f54aa95e33cbec2"
		TestPrivateKey3 = "923cebb3d8809d3caf09faa74ae2a39c23824a6fe75c44cab2a73dc6a0f3b606"
	)

	t.Run("should allowed only post request", func(t *testing.T) {
		// Arrange
		getRequest, _ := http.NewRequest(http.MethodGet, SendMailPath, nil)
		putRequest, _ := http.NewRequest(http.MethodPut, SendMailPath, nil)
		deleteRequest, _ := http.NewRequest(http.MethodDelete, SendMailPath, nil)
		postRequest, _ := http.NewRequest(http.MethodPost, SendMailPath, nil)

		// Act
		client := &http.Client{}
		getResponse, _ := client.Do(getRequest)
		putResponse, _ := client.Do(putRequest)
		deleteResponse, _ := client.Do(deleteRequest)
		postResponse, _ := client.Do(postRequest)

		// Assert
		assert.Equal(t, http.StatusMethodNotAllowed, getResponse.StatusCode)
		assert.Equal(t, http.StatusMethodNotAllowed, putResponse.StatusCode)
		assert.Equal(t, http.StatusMethodNotAllowed, deleteResponse.StatusCode)
		assert.NotEqual(t, http.StatusMethodNotAllowed, postResponse.StatusCode)
	})

	t.Run("should return unauthorized when request does not have public key in header", func(t *testing.T) {
		// Arrange
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, nil)

		// Act
		client := &http.Client{}
		response1, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response1.StatusCode)
	})

	t.Run("should return unauthorized when public key in header is invalid", func(t *testing.T) {
		// Arrange
		badPublicKeyHeader := "bad key heehee! ow!"
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, nil)
		request.Header.Add("x-public-key", badPublicKeyHeader)

		// Act
		client := &http.Client{}
		response1, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response1.StatusCode)
	})

	t.Run("should return bad request when user send bad request body", func(t *testing.T) {
		// Arrange
		sendAccount, newAccountErr1 := account.ConnectAccount(TestPrivateKey1)
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, nil)
		request.Header.Add("x-public-key", sendAccount.GetAddress())

		// Act
		client := &http.Client{}
		response1, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr1, newReqErr, sendReqErr)
		assert.Equal(t, http.StatusBadRequest, response1.StatusCode)
	})

	t.Run("should return unauthorized when request contains invalid signature", func(t *testing.T) {
		// Arrange
		sendAccount, newAccountErr1 := account.ConnectAccount(TestPrivateKey1)
		receivedAccount, newAccountErr2 := account.ConnectAccount(TestPrivateKey2)
		badAccount, newAccountErr3 := account.ConnectAccount(TestPrivateKey3)
		recipient := receivedAccount.GetAddress()
		subject := "test subject status created to " + recipient
		body := "test mail body to " + recipient
		message, newMsgErr := request.NewSendEmail(recipient, subject, body)
		signedBadMessage, signErr := badAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedBadMessage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, payLoad)
		request.Header.Add("x-public-key", sendAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)
		resultBytes, readErr := io.ReadAll(response.Body)
		result := model.SendMailResponse{}
		unmarshalErr := json.Unmarshal(resultBytes, &result)

		// Assert
		util.AssertNoAnyError(t, newAccountErr1, newAccountErr2, newAccountErr3, newMsgErr, signErr)
		util.AssertNoAnyError(t, marshalErr, newReqErr, sendReqErr, readErr)
		assert.Error(t, unmarshalErr)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("should return bad request when user send invalid recipient public key", func(t *testing.T) {
		// Arrange
		testAccount, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		badRecipientPublicKey := "bad key heehee! ow!"
		subject := "test subject"
		body := "test mail body"
		message, newMsgErr := request.NewSendEmail(badRecipientPublicKey, subject, body)
		signedMessage, signErr := testAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMessage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, payLoad)
		request.Header.Add("x-public-key", testAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)
		resultBytes, readErr := io.ReadAll(response.Body)
		result := model.SendMailResponse{}
		unmarshalErr := json.Unmarshal(resultBytes, &result)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, newMsgErr, signErr, marshalErr)
		util.AssertNoAnyError(t, newReqErr, sendReqErr, readErr)
		assert.Error(t, unmarshalErr)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return status created and mail id", func(t *testing.T) {
		// Arrange
		sendAccount, newAccountErr1 := account.ConnectAccount(TestPrivateKey1)
		receivedAccount, newAccountErr2 := account.ConnectAccount(TestPrivateKey2)
		recipient := receivedAccount.GetAddress()
		subject := "test subject status created to " + recipient
		body := "test mail body to " + recipient
		message, newMsgErr := request.NewSendEmail(recipient, subject, body)
		signedMessage, signErr := sendAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMessage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, payLoad)
		request.Header.Add("x-public-key", sendAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)
		resultBytes, readErr := io.ReadAll(response.Body)
		result := model.SendMailResponse{}
		unmarshalErr := json.Unmarshal(resultBytes, &result)

		// Assert
		util.AssertNoAnyError(t, marshalErr, newAccountErr1, newAccountErr2, newMsgErr, signErr)
		util.AssertNoAnyError(t, newReqErr, sendReqErr, readErr, unmarshalErr)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("recipient should receive mail", func(t *testing.T) {
		// Arrange: send mail
		sendAccount, _ := account.ConnectAccount(TestPrivateKey1)
		receivedAccount, _ := account.ConnectAccount(TestPrivateKey3)
		recipient := receivedAccount.GetAddress()
		subject := recipient + "should receive this subject"
		body := recipient + "should receive this body"
		sendMailMessage, _ := request.NewSendEmail(recipient, subject, body)
		signedMessage, _ := sendAccount.Sign(sendMailMessage)
		requestBody := model.RequestBody{
			Data:      string(sendMailMessage),
			Signature: signedMessage,
		}
		requestBodyByte, _ := json.Marshal(requestBody)
		payLoad := strings.NewReader(string(requestBodyByte))
		sendMailRequest, _ := http.NewRequest(http.MethodPost, SendMailPath, payLoad)
		sendMailRequest.Header.Add("x-public-key", sendAccount.GetAddress())
		client := &http.Client{}
		client.Do(sendMailRequest)

		// Act: read mail inbox
		readMailMessage, _ := request.NewGetInbox()
		signedReadMailMessage, _ := receivedAccount.Sign(readMailMessage)
		readMailRequestBody := model.RequestBody{
			Data:      string(readMailMessage),
			Signature: signedReadMailMessage,
		}
		readMailRequestBodyByte, _ := json.Marshal(readMailRequestBody)
		queryParams := "?page=1&limit=10"
		apiPath := InboxPath + queryParams
		readMailPayLoad := strings.NewReader(string(readMailRequestBodyByte))
		readMailRequest, _ := http.NewRequest(http.MethodPost, apiPath, readMailPayLoad)
		readMailRequest.Header.Add("x-public-key", receivedAccount.GetAddress())
		readMailResponse, _ := client.Do(readMailRequest)
		readMailBytes, _ := io.ReadAll(readMailResponse.Body)
		inbox := model.InboxResponse{}
		unmarshalErr := json.Unmarshal(readMailBytes, &inbox)

		// Assert
		assert.Equal(t, http.StatusOK, readMailResponse.StatusCode)
		assert.NoError(t, unmarshalErr)
		assert.Equal(t, 1, len(inbox.Inbox))
		assert.Equal(t, subject, inbox.Inbox[0].Subject)
		assert.Equal(t, body, inbox.Inbox[0].Body)
		assert.Equal(t, sendAccount.GetAddress(), inbox.Inbox[0].From)
	})

	t.Run("should return unauthorized when request uuid is duplicated", func(t *testing.T) {
		// Arrange
		sendAccount, newAccountErr1 := account.ConnectAccount(TestPrivateKey1)
		receivedAccount, newAccountErr2 := account.ConnectAccount(TestPrivateKey2)
		recipient := receivedAccount.GetAddress()
		subject := "test subject uuid duplicate to " + recipient
		body := "test mail body to " + recipient
		message, newMsgErr := request.NewSendEmail(recipient, subject, body)
		signedMessage, signErr := sendAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMessage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		payLoad1 := strings.NewReader(string(requestBodyByte))
		payLoad2 := strings.NewReader(string(requestBodyByte))
		request1, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, payLoad1)
		request1.Header.Add("x-public-key", sendAccount.GetAddress())
		request2, newReqErr2 := http.NewRequest(http.MethodPost, SendMailPath, payLoad2)
		request2.Header.Add("x-public-key", sendAccount.GetAddress())

		// Act
		client := &http.Client{}
		response1, sendReqErr1 := client.Do(request1)
		response2, sendReqErr2 := client.Do(request2) // duplicated uuid request

		// Assert
		util.AssertNoAnyError(t, newAccountErr1, newAccountErr2, newMsgErr, signErr)
		util.AssertNoAnyError(t, marshalErr, newReqErr, sendReqErr1, sendReqErr2, newReqErr2)
		assert.Equal(t, http.StatusCreated, response1.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, response2.StatusCode)
	})

	t.Run("should return unauthorized when request contains timeout timestamp", func(t *testing.T) {
		// Arrange
		last3minutes1second := time.Now().Add(-3 * time.Minute).Add(-1 * time.Second) // timeout is 3 minutes
		sendAccount, newAccountErr1 := account.ConnectAccount(TestPrivateKey1)
		receivedAccount, newAccountErr2 := account.ConnectAccount(TestPrivateKey2)
		recipient := receivedAccount.GetAddress()
		subject := "test subject timeout to " + recipient
		body := "test mail body to " + recipient
		sendEmail := request.SendEmailRequest{
			ID:        uuid.New(),
			Timestamp: last3minutes1second.Format(time.RFC3339),
			Recipient: recipient,
			Subject:   subject,
			Body:      body,
		}
		message, newMsgErr := json.Marshal(sendEmail)
		signedMessage, signErr := sendAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMessage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, SendMailPath, payLoad)
		request.Header.Add("x-public-key", sendAccount.GetAddress())
		request.Header.Add("x-timestamp", "test")

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr1, newAccountErr2, newMsgErr, signErr)
		util.AssertNoAnyError(t, marshalErr, newReqErr, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
}
