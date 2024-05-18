package tests

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInbox(t *testing.T) {

	t.Skip("make test faster")

	const (
		testPrivateKey = "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247"
		mainPath       = "../cmd/main.go"
		testQuery      = "../tests/util/query.test.json"
	)

	t.Run("should handle invalid query path", func(t *testing.T) {
		invalidPath := "../invalid/path/HeeHee/AAOW!.jackson"
		cmd := exec.Command("go", "run", mainPath, "-inbox", invalidPath, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "query file not found, invalid path or file name\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should handle invalid query json with value type", func(t *testing.T) {
		invalidQuery := "../tests/util/bad-query-01-type.test.json"
		cmd := exec.Command("go", "run", mainPath, "-inbox", invalidQuery, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		expectError := "invalid query json: please use this format\n\t{ \"page\":int, \"limit\":int }\n"
		assert.Equal(t, expectError, stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should handle invalid query json with missing field", func(t *testing.T) {
		invalidQuery := "../tests/util/bad-query-02-key.test.json"
		cmd := exec.Command("go", "run", mainPath, "-inbox", invalidQuery, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		expectError := "invalid query json: please use this format\n\t{ \"page\":int, \"limit\":int }\n"
		assert.Equal(t, expectError, stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should require user private key", func(t *testing.T) {
		cmd := exec.Command("go", "run", mainPath, "-inbox", testQuery)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "user credential is required\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should handle invalid user private key with wrong length", func(t *testing.T) {
		invalidCredential := "1baa694c49154f63b15" // wrong length
		cmd := exec.Command("go", "run", mainPath, "-inbox", testQuery, "-user", invalidCredential)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "invalid user credential: credential should be hex with 64 characters long\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should handle invalid user private key with wrong format", func(t *testing.T) {
		invalidCredential := "x#aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // non-hex
		cmd := exec.Command("go", "run", mainPath, "-inbox", testQuery, "-user", invalidCredential)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		assert.Error(t, err)
		assert.Equal(t, "invalid user credential: credential should be hex with 64 characters long\n", stdout.String())
		assert.Equal(t, "exit status 1\n", stderr.String())
	})

	t.Run("should send correct api request to server", func(t *testing.T) {
		cmd := exec.Command("go", "run", mainPath, "-inbox", testQuery, "-user", testPrivateKey)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		cmdErr := cmd.Run()

		result := stdout.String()
		haveInbox := strings.Contains(result, "inbox")
		haveTotal := strings.Contains(result, "total")

		assert.NoError(t, cmdErr)
		assert.True(t, haveInbox)
		assert.True(t, haveTotal)
		assert.Empty(t, stderr.String())
	})
}
