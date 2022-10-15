package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type Note struct {
	mgm.DefaultModel `bson:",inline"`
	UserId           string `bson:"userId"`
	CipherText       []byte `bson:"cipherText"`
	KeyVersion       int64
}

func (k Note) GetIdStr() string {
	return k.ID.Hex()
}

func (k Note) IsIdEmpty() bool {
	return k.ID.IsZero()
}

func (k *Note) CollectionName() string {
	return "note"
}

func (k Note) GetCreatedAt() time.Time {
	return k.CreatedAt
}

func (k Note) GetUpdatedAt() time.Time {
	return k.UpdatedAt
}
