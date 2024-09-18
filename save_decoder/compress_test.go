package save_decoder

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompress(t *testing.T) {
	// First, decrypt an existing known save
	jsonSave, err := Decrypt(v64_slot1)
	assert.NoError(t, err)

	//fmt.Println(string(jsonSave))

	// Then, compress it
	compressed, err := CompressSave(jsonSave)
	assert.NoError(t, err)

	fmt.Println(compressed)

	// Decompress it
	decompressed, err := DecompressSave(compressed)
	assert.NoError(t, err)

	//fmt.Println(string(decompressed))

	// Finally, compare the decompressed save with the original decrypted save
	assert.Equal(t, jsonSave, decompressed)
}
