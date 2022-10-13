package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type UserKeyGenerator struct {
	mgm.DefaultModel  `bson:",inline"`
	UserId            string `bson:"userId"`
	KeyDerivationSalt []byte `bson:"keyDerivationSalt"`
	KeyHash           []byte `bson:"keyHash"`
	KeyVersion        int64
}

func (k UserKeyGenerator) GetIdStr() string {
	return k.ID.Hex()
}

func (k UserKeyGenerator) IsIdEmpty() bool {
	return k.ID.IsZero()
}

func (k *UserKeyGenerator) CollectionName() string {
	return "userKeys"
}

func (k UserKeyGenerator) GetCreatedAt() time.Time {
	return k.CreatedAt
}

func (k UserKeyGenerator) GetUpdatedAt() time.Time {
	return k.UpdatedAt
}
