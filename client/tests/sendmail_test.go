package tests

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMail(t *testing.T) {

	const (
		testPrivateKey = "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247"
		mainPath       = "../cmd/main.go"
		testMail       = "../tests/util/good.kmail.json"
	)

	t.Run("should handle invalid kmail path", func(t *testing.T) {
		invalidPath := "../invalid/path/HeeHee/AAOW!.jackson"
		cmd := exec.Command("go", "run", mainPath, "-send", invalidPath, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "mail file not found, invalid path or file name\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should handle invalid kmail json with value type", func(t *testing.T) {
		badMail := "../tests/util/bad-mail-01-value.kmail.json"
		cmd := exec.Command("go", "run", mainPath, "-send", badMail, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "invalid kmail json: please use this format\n\t{ \"to\":string, \"subject\":string, \"body\":string }\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should handle invalid user private key", func(t *testing.T) {
		badPrivateKey := "y57g8glkrj09ug039[]"
		cmd := exec.Command("go", "run", mainPath, "-send", testMail, "-user", badPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "invalid user credential: credential should be hex with 64 characters long\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should send correct api request to server", func(t *testing.T) {
		// send mail correctly
		cmd := exec.Command("go", "run", mainPath, "-send", testMail, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		result := stdout.String()
		haveID := strings.Contains(result, "id")

		assert.NoError(t, err)
		assert.True(t, haveID)
		assert.Empty(t, stderr.String())
	})
}
