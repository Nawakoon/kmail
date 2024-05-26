package server_test

import (
	"encoding/json"
	"fmt"
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

func ReadMailTestCases(t *testing.T) {

	const (
		ReadMailPath    = "http://localhost:8080/mail"
		SendMailPath    = "http://localhost:8080/mail/send"
		TestPrivateKey1 = "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247"
		TestPrivateKey2 = "fd778940ddae63e19e5d2a05604a4d0eaec18b977801299a7f54aa95e33cbec2"
		TestPrivateKey3 = "923cebb3d8809d3caf09faa74ae2a39c23824a6fe75c44cab2a73dc6a0f3b606"
	)

	var mailUUID uuid.UUID

	beforeAll := func() {
		allMails0 := retrieveMails()
		fmt.Printf("allMails0: %v\n", allMails0)
		sendAccount, _ := account.ConnectAccount(TestPrivateKey1)
		receivedAccount, _ := account.ConnectAccount(TestPrivateKey2)
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
		response, err := client.Do(sendMailRequest)
		if err != nil {
			assert.Fail(t, "can not create mail for test")
			t.Skip("can not create mail for test")
		}
		responseBytes, _ := io.ReadAll(response.Body)
		mail := model.SendMailResponse{}
		err = json.Unmarshal(responseBytes, &mail)
		if err != nil {
			assert.Fail(t, "json unmarshal error")
			t.Skip("json unmarshal error")
		}

		mailUUID = mail.ID
	}

	beforeAll()

	t.Run("should allow only post request", func(t *testing.T) {
		// Arrange
		getRequest, _ := http.NewRequest(http.MethodGet, ReadMailPath, nil)
		putRequest, _ := http.NewRequest(http.MethodPut, ReadMailPath, nil)
		deleteRequest, _ := http.NewRequest(http.MethodDelete, ReadMailPath, nil)
		postRequest, _ := http.NewRequest(http.MethodPost, ReadMailPath, nil)

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

	t.Run("should return unauthorized when request does not have a public key in header", func(t *testing.T) {
		// Arrange
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, nil)

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("should return unauthorized when public key in header is invalid", func(t *testing.T) {
		// Arrange
		badPublicKeyHeader := "bad key heehee! ow!"
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, nil)
		request.Header.Add("x-public-key", badPublicKeyHeader)

		// Act
		client := &http.Client{}
		response1, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response1.StatusCode)
	})

	t.Run("should return bad request when user send invalid request uuid in request body", func(t *testing.T) {
		// Arrange
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		type BadMessage struct {
			ID        string    `json:"id"`
			Timestamp string    `json:"timestamp"`
			EmailID   uuid.UUID `json:"email_id"`
		}
		badMessage := BadMessage{
			ID:        "bad id hee hee",
			Timestamp: time.Now().Format(time.RFC3339),
			EmailID:   mailUUID,
		}
		badMessageByte, err := json.Marshal(badMessage)
		if err != nil {
			assert.Fail(t, "json marshal error")
			t.Skip("json marshal error")
		}
		signMessage, signMsgErr := account.Sign(badMessageByte)
		requestBody := model.RequestBody{
			Data:      string(badMessageByte),
			Signature: signMessage,
		}
		requestBodyByte, err := json.Marshal(requestBody)
		if err != nil {
			assert.Fail(t, "json marshal error")
			t.Skip("json marshal error")
		}
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, newReqErr, sendReqErr, signMsgErr)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return bad request when user send invalid timestamp in request body", func(t *testing.T) {
		// Arrange
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		badGetMailMessage := request.GetEmailRequest{
			ID:        uuid.New(),
			Timestamp: "bad timestamp",
			EmailID:   mailUUID,
		}
		badMessageByte, marshalMsgErr := json.Marshal(badGetMailMessage)
		signMessage, signMsgErr := account.Sign(badMessageByte)
		requestBody := model.RequestBody{
			Data:      string(badMessageByte),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, marshalMsgErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return bad request when user send invalid mail id in request body", func(t *testing.T) {
		// Arrange
		type BadMessage struct {
			ID        uuid.UUID `json:"id"`
			Timestamp string    `json:"timestamp"`
			EmailID   string    `json:"email_id"`
		}
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		badGetMailMessage := BadMessage{
			ID:        uuid.New(),
			Timestamp: time.Now().Format(time.RFC3339),
			EmailID:   "bad mail id hee hee",
		}
		badMessageByte, marshalMsgErr := json.Marshal(badGetMailMessage)
		signMessage, signMsgErr := account.Sign(badMessageByte)
		requestBody := model.RequestBody{
			Data:      string(badMessageByte),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, marshalMsgErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return unauthorized when request contains invalid signature", func(t *testing.T) {
		goodAccount, newGoodAccErr := account.ConnectAccount(TestPrivateKey1)
		badAccount, newBadAccErr := account.ConnectAccount(TestPrivateKey2)
		message, newMsgErr := request.NewGetEmail(mailUUID)
		messageByte, marshalMsgErr := json.Marshal(message)
		badSignMessage, signMsgErr := badAccount.Sign(messageByte)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: badSignMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", goodAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newGoodAccErr, newBadAccErr, newMsgErr, marshalMsgErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("should return not found when mail is not found in database", func(t *testing.T) {
		// Arrange
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		randomMailID := uuid.New()
		message, newMsgErr := request.NewGetEmail(randomMailID)
		signMessage, signMsgErr := account.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, newMsgErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("should return OK and mail data when sender account read the mail", func(t *testing.T) {
		// Arrange
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		message, newMsgErr := request.NewGetEmail(mailUUID)
		signMessage, signMsgErr := account.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)
		responseBytes, readErr := io.ReadAll(response.Body)
		mail := model.Mail{}
		unmarshalErr := json.Unmarshal(responseBytes, &mail)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, newMsgErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr, readErr, unmarshalErr)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, mailUUID, mail.ID)
	})

	t.Run("should return OK and mail data when recipient account read the mail", func(t *testing.T) {
		// Arrange
		recipientAccount, newAccountErr := account.ConnectAccount(TestPrivateKey2)
		message, newMsgErr := request.NewGetEmail(mailUUID)
		signMessage, signMsgErr := recipientAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", recipientAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)
		responseBytes, readErr := io.ReadAll(response.Body)
		mail := model.Mail{}
		unmarshalErr := json.Unmarshal(responseBytes, &mail)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, newMsgErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr, readErr, unmarshalErr)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, mailUUID, mail.ID)
	})

	t.Run("should return not found when account is not the recipient or sender of the mail", func(t *testing.T) {
		// Arrange
		unauthorizeAccount, newUnauthorizeAccErr := account.ConnectAccount(TestPrivateKey3)
		message, newMsgErr := request.NewGetEmail(mailUUID)
		unauthorizeSignMessage, unauthorizeSignMsgErr := unauthorizeAccount.Sign(message)
		unauthorizeRequestBody := model.RequestBody{
			Data:      string(message),
			Signature: unauthorizeSignMessage,
		}
		unauthorizeRequestBodyByte, marshalUnauthorizeByteErr := json.Marshal(unauthorizeRequestBody)
		unauthorizePayload := strings.NewReader(string(unauthorizeRequestBodyByte))
		unauthorizeRequest, newUnauthorizeReqErr := http.NewRequest(http.MethodPost, ReadMailPath, unauthorizePayload)
		unauthorizeRequest.Header.Add("x-public-key", unauthorizeAccount.GetAddress())

		// Act
		client := &http.Client{}
		unauthorizeResponse, sendUnauthorizeReqErr := client.Do(unauthorizeRequest)

		// Assert
		util.AssertNoAnyError(t,
			newUnauthorizeAccErr, newMsgErr, marshalUnauthorizeByteErr, newUnauthorizeReqErr,
			sendUnauthorizeReqErr, unauthorizeSignMsgErr,
		)
		assert.Equal(t, http.StatusNotFound, unauthorizeResponse.StatusCode)
	})

	t.Run("should return unauthorized when request uuid is duplicated", func(t *testing.T) {
		// Arrange
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		message, newMsgErr := request.NewGetEmail(mailUUID)
		signMessage, signMsgErr := account.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload1 := strings.NewReader(string(requestBodyByte))
		payload2 := strings.NewReader(string(requestBodyByte))
		request1, newReq1Err := http.NewRequest(http.MethodPost, ReadMailPath, payload1)
		request2, newReq2Err := http.NewRequest(http.MethodPost, ReadMailPath, payload2)
		request1.Header.Add("x-public-key", account.GetAddress())
		request2.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response1, sendReq1Err := client.Do(request1)
		time.Sleep(500 * time.Millisecond)
		response2, sendReq2Err := client.Do(request2)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, newMsgErr, marshalByteErr, newReq1Err, newReq2Err, sendReq1Err, sendReq2Err, signMsgErr)
		assert.Equal(t, http.StatusOK, response1.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, response2.StatusCode)
	})

	t.Run("should return unauthorized when request contains timeout timestamp", func(t *testing.T) {
		// Arrange
		account, newAccountErr := account.ConnectAccount(TestPrivateKey1)
		last3minutes1second := time.Now().Add(-3 * time.Minute).Add(-1 * time.Second)
		getEmail := request.GetEmailRequest{
			ID:        uuid.New(),
			Timestamp: last3minutes1second.Format(time.RFC3339),
			EmailID:   mailUUID,
		}
		message, jsonEncodeErr := json.Marshal(getEmail)
		signMessage, signMsgErr := account.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signMessage,
		}
		requestBodyByte, marshalByteErr := json.Marshal(requestBody)
		payload := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, ReadMailPath, payload)
		request.Header.Add("x-public-key", account.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		util.AssertNoAnyError(t, newAccountErr, jsonEncodeErr, marshalByteErr, newReqErr, sendReqErr, signMsgErr)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
}

func retrieveMails() []model.MailEntity {
	testDatabase, err := util.NewTestDatabase()
	if err != nil {
		return []model.MailEntity{}
	}

	var mails []model.MailEntity
	rows, err := testDatabase.DB.Query("SELECT * FROM mail")
	if err != nil {
		return []model.MailEntity{}
	}
	defer rows.Close()

	for rows.Next() {
		var mail model.MailEntity
		err := rows.Scan(&mail.ID, &mail.Recipient, &mail.Sender, &mail.MailSubject, &mail.Body, &mail.SentAt)
		if err != nil {
			return []model.MailEntity{}
		}
		mails = append(mails, mail)
	}

	return mails
}
