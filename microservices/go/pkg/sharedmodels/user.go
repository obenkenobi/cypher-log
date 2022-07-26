package sharedmodels

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	AuthId           string `json:"authId" bson:"authId"`
	UserId           string `json:"userId" bson:"userId"`
	UserName         string `json:"userName" bson:"userName"`
	DisplayName      string `json:"displayName" bson:"displayName"`
	UserCreatedAt    int64  `json:"userCreatedAt" bson:"userCreatedAt"`
	UserUpdatedAt    int64  `json:"userUpdatedAt" bson:"userUpdatedAt"`
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
