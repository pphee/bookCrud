package util

import (
	"testing"
)

func TestEncryptionService(t *testing.T) {
	validKey := []byte("your-encryption-key-here")
	invalidKey := []byte("short-key")

	service := NewEncryptionService(validKey)

	originalText := "Hello, world!"
	encryptedText, err := service.Encrypt(originalText)
	if err != nil {
		t.Errorf("Failed to encrypt: %s", err)
	}
	if encryptedText == originalText {
		t.Errorf("Encryption failed, encrypted text is the same as the original")
	}

	decryptedText, err := service.Decrypt(encryptedText)
	if err != nil {
		t.Errorf("Failed to decrypt: %s", err)
	}
	if decryptedText != originalText {
		t.Errorf("Decryption failed, decrypted text is not the same as the original")
	}

	_, err = service.Encrypt("")
	if err != nil {
		t.Errorf("Encrypting empty string should not fail: %s", err)
	}
	_, err = service.Decrypt("")
	if err == nil {
		t.Errorf("Decrypting empty string should fail")
	}

	invalidService := NewEncryptionService(invalidKey)
	_, err = invalidService.Encrypt(originalText)
	if err == nil {
		t.Errorf("Encryption with invalid key should fail")
	}
	_, err = invalidService.Decrypt(encryptedText)
	if err == nil {
		t.Errorf("Decryption with invalid key should fail")
	}

	_, err = service.Decrypt("not-base64-encoded-text")
	if err == nil {
		t.Errorf("Decrypting invalid ciphertext should fail")
	}
}
