package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	AuthId           string `json:"authId" bson:"authId"`
	UserName         string `json:"userName" bson:"userName"`
	DisplayName      string `json:"displayName" bson:"displayName"`
	Distributed      bool   `json:"distributed" bson:"distributed"`
	ToBeDeleted      bool   `json:"toBeDeleted" bson:"toBeDeleted"`
}

func (u User) GetIdStr() string {
	return u.ID.Hex()
}

func (u User) IsIdEmpty() bool {
	return u.ID.IsZero()
}

func (u *User) CollectionName() string {
	return "users"
}

func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u User) WillNotDeleted() bool {
	return !u.ToBeDeleted
}
