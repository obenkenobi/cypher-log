package models

// PrimaryAppSecretRef contains the kid of the main app secret used to initialize sessions
type PrimaryAppSecretRef struct {
	Kid string `json:"kid"`
}
