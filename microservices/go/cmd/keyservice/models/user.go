package models

import (
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
)

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
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

func (u User) IsChangedFromExternalDto(extUserDto userdtos.UserDto) bool {
	return extUserDto.Id != u.UserId ||
		extUserDto.UserName != u.UserName ||
		extUserDto.DisplayName != u.DisplayName ||
		extUserDto.CreatedAt != u.UserCreatedAt ||
		extUserDto.UpdatedAt != u.UserCreatedAt

}
