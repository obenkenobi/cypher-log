package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type KeyData struct {
	mgm.DefaultModel  `bson:",inline"`
	UserId            string `bson:"userId"`
	KeyDerivationSalt []byte `bson:"keyDerivationSalt"`
	KeyHash           []byte `bson:"keyHash"`
}

func (k KeyData) GetIdStr() string {
	return k.ID.Hex()
}

func (k KeyData) IsIdEmpty() bool {
	return k.ID.IsZero()
}

func (k *KeyData) CollectionName() string {
	return "users"
}

func (k KeyData) GetCreatedAt() time.Time {
	return k.CreatedAt
}

func (k KeyData) GetUpdatedAt() time.Time {
	return k.UpdatedAt
}
