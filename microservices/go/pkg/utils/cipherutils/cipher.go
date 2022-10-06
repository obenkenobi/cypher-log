package cipherutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/exp/slices"
)

const aesKeyLength = 32

// EncryptAES encrypts your data using AES
func EncryptAES(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, err := generateRandomBytes(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// DecryptAES decrypts your data using AES
func DecryptAES(key, data []byte) ([]byte, error) {
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

// GenerateRandomKeyAES generates a random 32 byte encryption key to be used with AES
func GenerateRandomKeyAES() ([]byte, error) {
	return generateRandomBytes(aesKeyLength)
}

// DeriveAESKeyFromPassword generates an encryption key with a password. Salt can be passed as
// nil if no salt is provided as a new salt will be returned alongside the new
// key. If a salt already exists, just pass the existing one, and it will be used
// to derive the new key.
func DeriveAESKeyFromPassword(password, salt []byte) (derivedKey []byte, passwordSalt []byte, err error) {
	if salt == nil {
		salt, err = generateRandomBytes(32)
		if err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 32768, 8, 1, aesKeyLength)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

// HashKeyBcrypt hashes an encryption key with bcrypts
func HashKeyBcrypt(key []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(key, 10)
}

// VerifyKeyHashBcrypt takes an existing key hash and key and verifies if the key
// is correct with bcrypt. If so, true is returned. Otherwise, false is returned.
// An error is returned if and only if the verification failed due to reasons
// other than the key not matching the hash.
func VerifyKeyHashBcrypt(keyHash []byte, key []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(keyHash, key)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, bcrypt.ErrMismatchedHashAndPassword
	}
	return true, nil
}

const sha256SaltLen = 32

// HashWithSaltSHA256 hashes the given value and returns a hash with an
// embedded salt
func HashWithSaltSHA256(value []byte) ([]byte, error) {
	hash, salt, err := hashValueSHA256(value, nil)
	return append(hash, salt...), err
}

// VerifyHashWithSaltSHA256 verifies the salted hash matches the value
func VerifyHashWithSaltSHA256(saltedHash, value []byte) (bool, error) {
	saltIndex := len(saltedHash) - sha256SaltLen
	hash := saltedHash[:saltIndex]
	salt := saltedHash[saltIndex:]
	hashToVerify, _, err := hashValueSHA256(value, salt)
	if err != nil {
		return false, err
	}
	return slices.Equal(hash, hashToVerify), nil
}

func hashValueSHA256(value, salt []byte) (hash []byte, hashSalt []byte, err error) {
	if salt == nil {
		salt, err = generateRandomBytes(sha256SaltLen)
		if err != nil {
			return nil, nil, err
		}
	}
	saltedValue := append(value, salt...)
	h := sha256.New()
	if _, err := h.Write(saltedValue); err != nil {
		return nil, nil, err
	}
	return h.Sum(nil), salt, nil
}

func generateRandomBytes(byteCount int) ([]byte, error) {
	bytes := make([]byte, byteCount)
	_, err := rand.Read(bytes)
	return bytes, err
}
