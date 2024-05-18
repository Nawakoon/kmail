package account_test

import (
	"passwordless-mail-client/pkg/account"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestPrivateKey = "1baa694c49154f63b1503c7138f184c80f221670f035403ff428a65183bab247"

func TestAccount(t *testing.T) {
	t.Run("should convert public key to hex string and convert it back", func(t *testing.T) {
		// Arrange
		testAccount, connectErr := account.ConnectAccount(TestPrivateKey)

		// Act
		accountAddress := account.PublicKeyToHex(testAccount.PublicKey)
		publicKey, covertParseErr := account.HexToPublicKey(accountAddress)

		// Assert
		assert.NoError(t, connectErr)
		assert.NoError(t, covertParseErr)
		assert.Equal(t, testAccount.PublicKey, publicKey)
	})
}
