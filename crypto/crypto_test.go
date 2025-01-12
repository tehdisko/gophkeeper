package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCryptoManager(t *testing.T) {
	password := "password"

	manager := NewCryptoManager(password)

	plaintext := "plaintext"

	encrypted, err := manager.Encrypt(plaintext)
	require.NoError(t, err)

	decrypted, err := manager.Decrypt(encrypted)
	require.NoError(t, err)

	require.Equal(t, plaintext, decrypted)
}
