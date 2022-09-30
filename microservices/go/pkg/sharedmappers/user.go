package sharedmappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded/embeddeduser"
)

func DistUserSaveDtoToUser(distUser userdtos.DistUserSaveDto, user *sharedmodels.User) {
	AuthIdAndUserPublicDtoToUserModel(distUser.AuthId, distUser.BaseUserPublicDto, user)
}

func AuthIdAndUserPublicDtoToUserModel(authId string, userDto embeddeduser.BaseUserPublicDto, user *sharedmodels.User) {
	user.AuthId = authId
	user.UserId = userDto.Id
	user.UserName = userDto.UserName
	user.DisplayName = userDto.DisplayName
	user.UserCreatedAt = userDto.CreatedAt
	user.UserUpdatedAt = userDto.UpdatedAt
}

func UserModelToUserBo(user sharedmodels.User, userBo *userbos.UserBo) {
	userBo.Id = user.UserId
	userBo.UserName = user.UserName
	userBo.DisplayName = user.DisplayName
	userBo.CreatedAt = user.UserCreatedAt
	userBo.UpdatedAt = user.UserUpdatedAt
}

func UserReadDtoAndIdentityToUserIdentityDto(
	userDto userdtos.UserReadDto,
	identity security.Identity,
	dest *userdtos.UserIdentityDto,
) {
	dest.AuthId = identity.GetAuthId()
	dest.Authorities = identity.GetAuthorities()
	dest.UserReadDto = userDto
}

func UserDtoAndIdentityToDistUserSaveDto(
	userDto userdtos.UserReadDto,
	identity security.Identity,
	dest *userdtos.DistUserSaveDto,
) {
	dest.AuthId = identity.GetAuthId()
	dest.BaseUserPublicDto = userDto.BaseUserPublicDto
}
