package mappers

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
)

func MapUserSaveDtoToUser(source *userdtos.UserSaveDto, dest *models.User) {
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
}

func MapUserToUserDto(source *models.User, dest *userdtos.UserDto) {
	if source == nil {
		*dest = userdtos.UserDto{}
	}
	dest.Id = source.ID.String()
	dest.Exists = false
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
}
