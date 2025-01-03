package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func Encrypt(plaintext string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	padding := block.BlockSize() - len(plaintext)%block.BlockSize()
	padText := append([]byte(plaintext), bytes.Repeat([]byte{byte(padding)}, padding)...)

	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(padText))
	mode.CryptBlocks(ciphertext, padText)

	result := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func Decrypt(ciphertextBase64 string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < block.BlockSize() {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:block.BlockSize()]
	ciphertext = ciphertext[block.BlockSize():]

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	padding := plaintext[len(plaintext)-1]
	if int(padding) > len(plaintext) {
		return "", errors.New("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-int(padding)]

	return string(plaintext), nil
}
