package crypt

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	t.Run("Encrypt and decrypt", func(t *testing.T) {
		content := []byte("This is the payload")
		key := "super_secure_testing_key"

		enrypted, err := Encrypt(content, key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		decrypted, err := Decrypt(enrypted, key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if string(content) != string(decrypted) {
			t.Fatalf("Expected decrypted text to be %s, got %s", content, decrypted)
		}
	})

	t.Run("Encrypt and decrypt string", func(t *testing.T) {
		content := "This is the payload"
		key := "super_secure_testing_key"

		enrypted, err := EncryptString(content, key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		decrypted, err := DecryptString(enrypted, key)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if content != decrypted {
			t.Fatalf("Expected decrypted text to be %s, got %s", content, decrypted)
		}
	})
}
