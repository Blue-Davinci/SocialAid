package data

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	KeyLength16 = 16 // 128 bits
	KeyLength24 = 24 // 192 bits
	KeyLength32 = 32 // 256 bits
)

var (
	ErrInvalidEncryptionKeyLength = errors.New("invalid key length, must be 16, 24, or 32 bytes")
)

// contextGenerator is a helper function that generates a new context.Context from a
// context.Context and a timeout duration. This is useful for creating new contexts with
// deadlines for outgoing requests in our data layer.
func contextGenerator(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

/**
====================================================================
Encryption Functions
====================================================================
**/
// generateSecurityKey() generates a cryptographically secure AES key of the given length
func generateSecurityKey(keyLength int) (string, error) {
	if keyLength != KeyLength16 && keyLength != KeyLength24 && keyLength != KeyLength32 {
		return "", ErrInvalidEncryptionKeyLength
	}
	// Generate a random key of the given length
	key := make([]byte, keyLength)
	// Read random bytes into the key slice
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	// Convert the key to a hex string
	hexKey := hex.EncodeToString(key)
	return hexKey, nil
}

// DecodeEncryptionKey decodes the encryption key from a hex string to a byte slice
func DecodeEncryptionKey(encryptionkey string) ([]byte, error) {
	if encryptionkey == "" {
		return nil, fmt.Errorf("encryption key cannot be empty")
	}
	decodedKey, err := hex.DecodeString(encryptionkey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}
	return decodedKey, nil
}

// EncryptData() encrypts the given data using the provided key and returns the encrypted data as a base64 encoded string
// The key must be a 16, 24, or 32 byte slice to select AES-128, AES-192, or AES-256 encryption
func EncryptData(data string, key []byte) (string, error) {
	// Create a new AES cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	// Switch the AES to Galois Counter Mode (GCM) for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	// Generate a random nonce for the encryption
	nonce := make([]byte, gcm.NonceSize())
	// Read random bytes into the nonce slice
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	// Encrypt the data using the GCM
	encrypted := gcm.Seal(nonce, nonce, []byte(data), nil)
	// Return the encrypted data as a base64 encoded string for storage
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// DecryptData() decrypts the given encrypted data using the provided key and returns the decrypted data as a string
// The key must be a 16, 24, or 32 byte slice to select AES-128, AES-192, or AES-256 encryption
func DecryptData(encryptedData string, key []byte) (string, error) {
	// Check if the encrypted data is empty
	if encryptedData == "" {
		return "", fmt.Errorf("encrypted data cannot be empty")
	}
	// Decode the base64 encoded encrypted data
	data, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	// Create a new AES cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	// Create a new GCM block from the AES cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	// Extract the nonce size from the GCM
	nonceSize := gcm.NonceSize()
	// Check if the encrypted data is too short
	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid encrypted data: insufficient length")
	}
	// Extract the nonce and ciphertext from the encrypted data
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	// Decrypt the data using the GCM
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}
	// Return the decrypted data as a string
	return string(decrypted), nil
}
