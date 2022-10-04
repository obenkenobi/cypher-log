package cipherutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const KeyLength = 32

// Encrypt encrypts your data
func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt decrypts your data
func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateRandomKey generates a random 32 byte encryption key
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, KeyLength)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// DeriveKey generates an encryption key with a password. Salt can be passed as
// nil if no salt is provided as a new salt will be returned alongside the new
// key. If a salt already exists, just pass the existing one, and it will be used
// to derive the new key.
func DeriveKey(password, salt []byte) (derivedKey []byte, passwordSalt []byte, err error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 32768, 8, 1, KeyLength)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

// HashKey hashes an encryption key
func HashKey(key []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(key, 10)
}

// VerifyKeyHash takes an existing key hash and key and verifies if the key is
// correct. If so, true is returned. Otherwise, false is returned. An error is
// returned if and only if the verification failed due to reasons other than the
// key not matching the hash.
func VerifyKeyHash(keyHash []byte, key []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(keyHash, key)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, bcrypt.ErrMismatchedHashAndPassword
	}
	return true, nil
}
