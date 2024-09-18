package save_decoder

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompress(t *testing.T) {
	// First, decrypt an existing known save
	jsonSave, err := Decrypt(EncryptedSave{v64_slot1})
	assert.NoError(t, err)

	//fmt.Println(string(jsonSave))

	// Then, compress it
	compressedSave, err := CompressSave(jsonSave)
	assert.NoError(t, err)

	fmt.Println(compressedSave)

	// Decompress it
	jsonSave2, err := DecompressSave(compressedSave)
	assert.NoError(t, err)

	//fmt.Println(string(decompressed))

	// Finally, compare the decompressed save with the original decrypted save
	assert.Equal(t, jsonSave.Data, jsonSave2.Data)
}
