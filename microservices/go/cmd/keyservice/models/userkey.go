package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type UserKey struct {
	mgm.DefaultModel  `bson:",inline"`
	UserId            string `bson:"userId"`
	KeyDerivationSalt []byte `bson:"keyDerivationSalt"`
	KeyHash           []byte `bson:"keyHash"`
}

func (k UserKey) GetIdStr() string {
	return k.ID.Hex()
}

func (k UserKey) IsIdEmpty() bool {
	return k.ID.IsZero()
}

func (k *UserKey) CollectionName() string {
	return "userKeys"
}

func (k UserKey) GetCreatedAt() time.Time {
	return k.CreatedAt
}

func (k UserKey) GetUpdatedAt() time.Time {
	return k.UpdatedAt
}
