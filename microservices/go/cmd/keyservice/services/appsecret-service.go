package services

import "github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"

// Todo: implement AppSecretService

type AppSecretService interface {
	GetAppSecret(kid string) models.AppSecret // Gets the app secret marked by the KID
	GetPrimaryAppSecret() models.AppSecret    // Get the main app secret
	RotatePrimaryAppSecret() models.AppSecret // Sets the n
}
