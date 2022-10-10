package models

type UserKeySession struct {
	// The user's encryption key that is encrypted with the app's secret
	EncryptedKey []byte
	AppSecretKid string
	TokenHash    []byte
}
