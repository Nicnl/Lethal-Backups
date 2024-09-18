package save_decoder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// First, decrypt an existing known save
	jsonSave1, err := Decrypt(EncryptedSave{v64_slot1})
	assert.NoError(t, err)

	// Then, encrypt it back
	encryptedSave, err := Encrypt(jsonSave1)
	assert.NoError(t, err)

	// Finally, decrypt the encrypted save and compare it with the original decrypted save
	jsonSave2, err := Decrypt(encryptedSave)
	assert.NoError(t, err)

	assert.Equal(t, jsonSave1, jsonSave2)
}
