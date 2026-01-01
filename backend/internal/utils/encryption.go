package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"mini-app-backend/internal/config"
)

type EncryptionUtil struct {
	secretKey []byte
}

func NewEncryptionUtil() *EncryptionUtil {
	config := config.Load()
	
	return &EncryptionUtil{
		secretKey: []byte(config.CookieEncryptionKey),
	}
}

func (e *EncryptionUtil) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (e *EncryptionUtil) Decrypt(ciphertext string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}