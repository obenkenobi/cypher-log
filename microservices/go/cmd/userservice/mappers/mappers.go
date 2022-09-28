package mappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
)

func MapUserSaveDtoToUser(source userdtos.UserSaveDto, dest *models.User) {
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
}

func MapUserToUserDto(source models.User, dest *userdtos.UserDto) {
	dest.Id = source.GetIdStr()
	dest.Exists = !source.IsIdEmpty()
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
	dest.CreatedAt = source.CreatedAt.UnixMilli()
	dest.UpdatedAt = source.UpdatedAt.UnixMilli()
}

func MapToUserDtoAndIdentityToUserIdentityDto(
	userDto userdtos.UserDto,
	identity security.Identity,
	dest *userdtos.UserIdentityDto,
) {
	*dest = userdtos.UserIdentityDto{
		AuthId:      identity.GetAuthId(),
		Authorities: identity.GetAuthorities(),
		User:        userDto,
	}
}

func MapToUserDtoAndIdentityToDistUserDto(
	userDto userdtos.UserDto,
	identity security.Identity,
	dest *userdtos.DistributedUserDto,
) {
	*dest = userdtos.DistributedUserDto{
		AuthId: identity.GetAuthId(),
		User:   userDto,
	}
}
