package save_decoder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

func Encrypt(jsonSave []byte) ([]byte, error) {
	// Generate random IV (Initialization Vector)
	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	// Derive the key using PBKDF2 with SHA1 hash algorithm
	key := pbkdf2.Key([]byte(PASSWORD), iv, ITERATIONS, DK_LEN, sha1.New)

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Use CBC mode
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad the data (PKCS7 padding)
	paddedData := PKCS7Pad(jsonSave, aes.BlockSize)

	// Encrypt the data
	encrypted := make([]byte, len(paddedData))
	mode.CryptBlocks(encrypted, paddedData)

	// Prepend IV to the encrypted data
	encrypted = append(iv, encrypted...)

	return encrypted, nil
}

// PKCS7Pad adds padding according to the PKCS7 standard.
func PKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
