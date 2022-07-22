package models

import "github.com/kamva/mgm/v3"

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	ProviderUserId   string `json:"providerUserId" bson:"providerUserId"`
	UserName         string `json:"userName" bson:"userName"`
	DisplayName      string `json:"displayName" bson:"displayName"`
}
