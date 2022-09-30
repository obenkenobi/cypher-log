package mappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
)

func UserSaveDtoToUser(source userdtos.UserSaveDto, dest *models.User) {
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
}

func UserToUserDto(source models.User, dest *userdtos.UserReadDto) {
	dest.Id = source.GetIdStr()
	dest.Exists = !source.IsIdEmpty()
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
	dest.CreatedAt = source.CreatedAt.UnixMilli()
	dest.UpdatedAt = source.UpdatedAt.UnixMilli()
}

func UserToDistUserDeleteDto(source models.User, dest *userdtos.DistUserDeleteDto) {
	dest.Id = source.GetIdStr()
}
