package server_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"passwordless-mail-server/pkg/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	StartServerTime = 1000 * time.Millisecond
	BaseApiPath     = "http://localhost:8080"
)

func ClearPort() {
	fmt.Println("clearing port...")

	cmd := exec.Command("bash", "../scripts/clear-port.sh")
	err := cmd.Run()

	if err != nil {
		fmt.Println("warning: can not clear port correctly")
	}

	fmt.Println("port cleared")
}

func StartServer() {
	fmt.Println("starting server...")

	scriptPath := "../scripts/start-server.sh"
	cmd := exec.Command("bash", scriptPath)

	err := cmd.Run()
	if err != nil {
		log.Fatal("error: can not start server")
	}

	time.Sleep(StartServerTime)

	fmt.Println("server started")
}

func TestServer(t *testing.T) {

	testDatabase, _ := util.NewTestDatabase()
	ClearPort()
	StartServer()
	defer testDatabase.DeleteItemsFromTable("mail")
	defer testDatabase.DeleteItemsFromTable("used_uuid")

	t.Run("should have healthy status", func(t *testing.T) {
		// Arrange
		apiPath := BaseApiPath + "/health"

		// Act
		resp, apiErr := http.Get(apiPath)
		body, readErr := io.ReadAll(resp.Body)

		// Assert
		assert.NoError(t, apiErr)
		assert.NoError(t, readErr)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "OK\n", string(body))
	})

	// t.Run("should handle /mail/inbox", InboxTestCases)

	t.Run("should handle /mail/send", SendMailTestCases)
}
