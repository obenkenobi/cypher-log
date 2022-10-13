package mappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded/embeddeduser"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
)

func UserSaveDtoToUser(source userdtos.UserSaveDto, dest *models.User) {
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
}

func UserToUserDto(source models.User, dest *userdtos.UserReadDto) {
	dest.Exists = !source.IsIdEmpty()
	userToUserPublicDto(source, &dest.BaseUserPublicDto)
}

func UserToUserChangeEventDto(source models.User, dest *userdtos.UserChangeEventDto) {
	dest.AuthId = source.AuthId
	userToUserPublicDto(source, &dest.BaseUserPublicDto)
}

func userToUserPublicDto(source models.User, dest *embeddeduser.BaseUserPublicDto) {
	sharedmappers.MapMongoModelToBaseCrudObject(&source, &dest.BaseCRUDObject)
	dest.UserName = source.UserName
	dest.DisplayName = source.DisplayName
}
