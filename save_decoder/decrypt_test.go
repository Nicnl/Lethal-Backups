package save_decoder

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	//go:embed save_files/v64/LCSaveFile1
	v64_slot1 []byte

	//go:embed save_files/v64/LCSaveFile2
	v64_slot2 []byte

	//go:embed save_files/v64/LCSaveFile3
	v64_slot3 []byte
)

func TestDecrypt(t *testing.T) {
	output, err := Decrypt(v64_slot1)
	assert.NoError(t, err)

	fmt.Println(string(output))
}
