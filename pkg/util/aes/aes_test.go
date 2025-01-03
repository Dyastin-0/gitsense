package aes

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	code := m.Run()

	os.Exit(code)
}

func TestEncryptDecryptPrivateKey(t *testing.T) {
	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	originalPrivateKey := os.Getenv("PRIVATE_KEY")

	encryptedPrivateKey, err := Encrypt(originalPrivateKey, encryptionKey)
	if err != nil {
		t.Fatalf("Failed to encrypt private key: %v", err)
	}

	decryptedPrivateKey, err := Decrypt(encryptedPrivateKey, encryptionKey)
	if err != nil {
		t.Fatalf("Failed to decrypt private key: %v", err)
	}

	if decryptedPrivateKey != originalPrivateKey {
		t.Errorf("Decrypted private key doesn't match original: expected %v, got %v", originalPrivateKey, decryptedPrivateKey)
	}
}
