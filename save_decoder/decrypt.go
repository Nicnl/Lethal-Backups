package save_decoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

// Parameters for Lethal Company
const (
	DK_LEN     = 16
	ITERATIONS = 100
	PASSWORD   = "lcslime14a5"
)

// Inspired from:
// https://github.com/squee72564/LC_Simple_Save_Editor/blob/c9782c9639eab0493b1a0a3111d4adab0cf40282/src/encryption/encryptTools.py

// PKCS7Unpad removes padding according to the PKCS7 standard.
func PKCS7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid PKCS7 padding (empty data)")
	}
	if len(data)%blockSize != 0 {
		return nil, errors.New("invalid PKCS7 padding (incorrect block size)")
	}

	paddingLength := int(data[len(data)-1])
	if paddingLength == 0 || paddingLength > blockSize {
		return nil, errors.New("invalid PKCS7 padding (padding size out of range)")
	}

	// Validate padding bytes
	for _, padByte := range data[len(data)-paddingLength:] {
		if int(padByte) != paddingLength {
			return nil, errors.New("invalid PKCS7 padding (incorrect padding bytes)")
		}
	}

	return data[:len(data)-paddingLength], nil
}

func Decrypt(encoded []byte) ([]byte, error) {
	// Extract IV (Initialization Vector) from the first 16 bytes of data
	iv := encoded[:16]
	encryptedData := encoded[16:]

	// Derive the key using PBKDF2 with SHA1 hash algorithm
	key := pbkdf2.Key([]byte(PASSWORD), iv, ITERATIONS, DK_LEN, sha1.New)

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Use CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data
	decrypted := make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	// Unpad the decrypted data (PKCS7 padding)
	decrypted, err = PKCS7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("failed to unpad decrypted data: %w", err)
	}

	// Convert decrypted data to string (assuming it was a UTF-8 encoded string)
	return decrypted, nil
}
