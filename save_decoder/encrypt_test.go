package save_decoder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// First, decrypt an existing known save
	decoded1, err := Decrypt(v64_slot1)
	assert.NoError(t, err)

	// Then, encrypt it back
	encoded, err := Encrypt(decoded1)
	assert.NoError(t, err)

	// Finally, decrypt the encrypted save and compare it with the original decrypted save
	decoded2, err := Decrypt(encoded)
	assert.NoError(t, err)

	assert.Equal(t, decoded1, decoded2)
}
