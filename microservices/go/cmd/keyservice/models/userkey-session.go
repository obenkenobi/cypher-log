package models

type UserKeySession struct {
	// The user's encryption key that is encrypted with the app's secret
	EncryptedKey []byte `json:"encryptedKey"`
	AppSecretKid string `json:"appSecretKid"`
	TokenHash    []byte `json:"tokenHash"`
}
