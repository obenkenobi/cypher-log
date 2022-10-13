package models

type UserKeySession struct {
	// The user's encryption key that is encrypted with the app's secret
	KeyCipher        []byte `json:"keyCipher"`
	UserIdCipher     []byte `json:"userIdCipher"`
	KeyVersionCipher []byte `json:"keyVersionCipher"`
	AppSecretKid     string `json:"appSecretKid"`
	TokenHash        []byte `json:"tokenHash"`
}
