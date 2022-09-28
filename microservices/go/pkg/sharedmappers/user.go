package sharedmappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedbusinessobjects"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
)

func MapDistUserToUser(distUser userdtos.DistributedUserDto, user *sharedmodels.User) {
	MapAuthIdAndUserDtoToUser(distUser.AuthId, distUser.User, user)
}

func MapAuthIdAndUserDtoToUser(authId string, userDto userdtos.UserDto, user *sharedmodels.User) {
	user.AuthId = authId
	user.UserId = userDto.Id
	user.UserName = userDto.UserName
	user.DisplayName = userDto.DisplayName
	user.UserCreatedAt = userDto.CreatedAt
	user.UserUpdatedAt = userDto.UpdatedAt
}

func MapUserToUserBo(user sharedmodels.User, userBo *sharedbusinessobjects.UserBo) {
	userBo.UserId = user.UserId
	userBo.UserName = user.UserName
	userBo.DisplayName = user.DisplayName
	userBo.UserCreatedAt = user.UserCreatedAt
	userBo.UserUpdatedAt = user.UserUpdatedAt
}
