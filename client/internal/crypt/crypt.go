package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func Encrypt(value []byte, key string) ([]byte, error) {
	aesBlock, err := aes.NewCipher(hashKey(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	cipheredText := gcm.Seal(nonce, nonce, value, nil)
	return cipheredText, nil
}

func EncryptString(value string, key string) (string, error) {
	data, err := Encrypt([]byte(value), key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func Decrypt(encrypted []byte, key string) ([]byte, error) {
	hashedKey := hashKey(key)
	aesBlock, err := aes.NewCipher(hashedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipheredText := encrypted[:nonceSize], encrypted[nonceSize:]

	return gcm.Open(nil, nonce, cipheredText, nil)
}

func DecryptString(encrypted string, key string) (string, error) {
	decoded, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	data, err := Decrypt(decoded, key)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func hashKey(key string) []byte {
	text := []byte(key)
	hash := sha256.Sum256(text)
	return hash[:]
}
