package mappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/bos"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
)

func MapDistUserToUser(distUser userdtos.DistributedUserDto, user *models.User) {
	MapAuthIdAndUserDtoToUser(distUser.AuthId, distUser.User, user)
}

func MapAuthIdAndUserDtoToUser(authId string, userDto userdtos.UserDto, user *models.User) {
	user.AuthId = authId
	user.UserId = userDto.Id
	user.UserName = userDto.UserName
	user.DisplayName = userDto.DisplayName
	user.UserCreatedAt = userDto.CreatedAt
	user.UserUpdatedAt = userDto.UpdatedAt
}

func MapUserToUserBo(user models.User, userBo *bos.UserBo) {
	userBo.UserId = user.UserId
	userBo.UserName = user.UserName
	userBo.DisplayName = user.DisplayName
	userBo.UserCreatedAt = user.UserCreatedAt
	userBo.UserUpdatedAt = user.UserUpdatedAt
}
