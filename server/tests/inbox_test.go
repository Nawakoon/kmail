package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"passwordless-mail-client/pkg/account"
	"passwordless-mail-client/pkg/request"
	"passwordless-mail-server/pkg/model"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func InboxTestCases(t *testing.T) {

	const (
		BaseInboxPath  = "http://localhost:8080/mail/inbox"
		TestPrivateKey = "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247"
	)

	t.Run("should allowed only post request", func(t *testing.T) {
		// Arrange
		getRequest, _ := http.NewRequest(http.MethodGet, BaseInboxPath, nil)
		putRequest, _ := http.NewRequest(http.MethodPut, BaseInboxPath, nil)
		deleteRequest, _ := http.NewRequest(http.MethodDelete, BaseInboxPath, nil)
		postRequest, _ := http.NewRequest(http.MethodPost, BaseInboxPath, nil)

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
		request, newReqErr := http.NewRequest(http.MethodPost, BaseInboxPath, nil)

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
		request, newReqErr := http.NewRequest(http.MethodPost, BaseInboxPath, nil)
		request.Header.Add("x-public-key", "test")

		// Act
		client := &http.Client{}
		response1, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusBadRequest, response1.StatusCode)
	})

	t.Run("should return unauthorized when request contains invalid signature", func(t *testing.T) {
		// Arrange
		requestBody := model.RequestBody{
			Data:      "test",
			Signature: []byte("test"),
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, BaseInboxPath, payLoad)
		request.Header.Add("x-public-key", "test")

		// Act
		client := &http.Client{}
		response1, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, marshalErr)
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response1.StatusCode)
	})

	t.Run("should return bad request when request contains invalid query params", func(t *testing.T) {
		// Arrange
		testAccount, _ := account.ConnectAccount(TestPrivateKey)
		message, _ := request.NewGetInbox()
		signedMassage, _ := testAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMassage,
		}
		requestBodyByte, _ := json.Marshal(requestBody)
		badQueryParams := "?page=heehee&limit=ow"
		apiPath := BaseInboxPath + badQueryParams
		payLoad := strings.NewReader(string(requestBodyByte))
		request, _ := http.NewRequest(http.MethodPost, apiPath, payLoad)
		request.Header.Add("x-public-key", testAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, _ := client.Do(request)

		// Assert
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("should return ok and mail inbox when user send request correctly", func(t *testing.T) {
		// Arrange
		testAccount, connectErr := account.ConnectAccount(TestPrivateKey)
		message, newMsgErr := request.NewGetInbox()
		signedMassage, signErr := testAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMassage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		queryParams := "?page=1&limit=10"
		apiPath := BaseInboxPath + queryParams
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, apiPath, payLoad)
		request.Header.Add("x-public-key", testAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)
		inboxBytes, readErr := io.ReadAll(response.Body)
		inbox := model.InboxResponse{}
		unmarshalErr := json.Unmarshal(inboxBytes, &inbox)

		// Assert
		assert.NoError(t, connectErr)
		assert.NoError(t, newMsgErr)
		assert.NoError(t, signErr)
		assert.NoError(t, marshalErr)
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.NoError(t, readErr)
		assert.NoError(t, unmarshalErr)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("should return unauthorize when request uuid is duplicate", func(t *testing.T) {
		// Arrange
		testAccount, connectErr := account.ConnectAccount(TestPrivateKey)
		message, newMsgErr := request.NewGetInbox()
		signedMassage, signErr := testAccount.Sign(message)
		requestBody := model.RequestBody{
			Data:      string(message),
			Signature: signedMassage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		queryParams := "?page=1&limit=10"
		apiPath := BaseInboxPath + queryParams
		payLoad1 := strings.NewReader(string(requestBodyByte))
		payLoad2 := strings.NewReader(string(requestBodyByte))
		request1, newReqErr1 := http.NewRequest(http.MethodPost, apiPath, payLoad1)
		request1.Header.Add("x-public-key", testAccount.GetAddress())
		request2, newReqErr2 := http.NewRequest(http.MethodPost, apiPath, payLoad2)
		request2.Header.Add("x-public-key", testAccount.GetAddress())

		// Act
		client := &http.Client{}
		response1, sendReq1Err := client.Do(request1)
		response2, sendReq2Err := client.Do(request2)

		// Assert
		assert.NoError(t, connectErr)
		assert.NoError(t, newMsgErr)
		assert.NoError(t, signErr)
		assert.NoError(t, marshalErr)
		assert.NoError(t, newReqErr1)
		assert.NoError(t, newReqErr2)
		assert.NoError(t, sendReq1Err)
		assert.NoError(t, sendReq2Err)
		assert.Equal(t, http.StatusOK, response1.StatusCode)
		assert.Equal(t, http.StatusUnauthorized, response2.StatusCode)
	})

	t.Run("should return unauthorize when request contains timeout timestamp", func(t *testing.T) {
		// Arrange
		testAccount, connectErr := account.ConnectAccount(TestPrivateKey)
		last3minutes1second := time.Now().Add(-3 * time.Minute).Add(-1 * time.Second)
		getInbox := request.GetInboxRequest{
			ID:        uuid.New(),
			Timestamp: last3minutes1second.Format(time.RFC3339),
		}
		inbox, jsonEncodeErr := json.Marshal(getInbox)

		signedMassage, signErr := testAccount.Sign(inbox)
		requestBody := model.RequestBody{
			Data:      string(inbox),
			Signature: signedMassage,
		}
		requestBodyByte, marshalErr := json.Marshal(requestBody)
		queryParams := "?page=1&limit=10"
		apiPath := BaseInboxPath + queryParams
		payLoad := strings.NewReader(string(requestBodyByte))
		request, newReqErr := http.NewRequest(http.MethodPost, apiPath, payLoad)
		request.Header.Add("x-public-key", testAccount.GetAddress())

		// Act
		client := &http.Client{}
		response, sendReqErr := client.Do(request)

		// Assert
		assert.NoError(t, connectErr)
		assert.NoError(t, jsonEncodeErr)
		assert.NoError(t, signErr)
		assert.NoError(t, marshalErr)
		assert.NoError(t, newReqErr)
		assert.NoError(t, sendReqErr)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})
}
