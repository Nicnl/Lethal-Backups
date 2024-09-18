package save_decoder

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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
	curSave, err := os.ReadFile("C:\\Users\\Nicnl\\AppData\\LocalLow\\ZeekerssRBLX\\Lethal Company\\LCSaveFile3")
	assert.NoError(t, err)

	output, err := Decrypt(curSave)
	assert.NoError(t, err)

	fmt.Println(string(output))
	os.WriteFile("C:\\Users\\Nicnl\\GolandProjects\\lethal_company_save_manager\\save_decoder\\out.json", []byte(output), 0644)
}
