package services

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/mappers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/streamx/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	AddUser(identity security.Identity, userSaveDto userdtos.UserSaveDto) single.Single[userdtos.UserDto]
	UpdateUser(identity security.Identity, userSaveDto userdtos.UserSaveDto) single.Single[userdtos.UserDto]
	GetByAuthId(tokenId string) single.Single[userdtos.UserDto]
	GetUserIdentity(identity security.Identity) single.Single[userdtos.UserIdentityDto]
}

type userServiceImpl struct {
	dbHandler      database.DBHandler
	userRepository repositories.UserRepository
	userBr         businessrules.UserBr
	errorService   apperrors.ErrorService
}

func (u userServiceImpl) AddUser(
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserDto] {
	userCreateValidationSrc := u.userBr.ValidateUserCreate(u.dbHandler.GetCtx(), identity, userSaveDto)
	userCreateSrc := single.FlatMap(userCreateValidationSrc, func([]apperrors.RuleError) single.Single[models.User] {
		user := models.User{}
		mappers.MapUserSaveDtoToUser(userSaveDto, &user)
		user.AuthId = identity.GetAuthId()
		return single.MapDerefPtr(u.userRepository.Create(u.dbHandler.GetCtx(), &user))
	})
	return single.Map(userCreateSrc, func(user models.User) userdtos.UserDto {
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		log.Debug("Created user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) UpdateUser(
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserDto] {
	userSearchSrc := u.userRepository.FindByAuthId(u.dbHandler.GetCtx(), identity.GetAuthId())
	userExistsSrc := single.MapWithError(
		userSearchSrc,
		func(userMaybe option.Maybe[models.User]) (models.User, error) {
			if user, ok := userMaybe.Get(); ok {
				return user, nil
			} else {
				err := apperrors.NewBadReqErrorFromRuleError(
					u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqItemsNotFound))
				return user, err
			}
		},
	)
	userValidatedSrc := single.FlatMap(userExistsSrc, func(existingUser models.User) single.Single[models.User] {
		validationSrc := u.userBr.ValidateUserUpdate(u.dbHandler.GetCtx(), userSaveDto, existingUser)
		return single.Map(validationSrc, func([]apperrors.RuleError) models.User { return existingUser })
	})
	userSavedSrc := single.FlatMap(userValidatedSrc, func(user models.User) single.Single[models.User] {
		mappers.MapUserSaveDtoToUser(userSaveDto, &user)
		return single.MapDerefPtr(u.userRepository.Update(u.dbHandler.GetCtx(), &user))
	})
	return single.Map(userSavedSrc, func(user models.User) userdtos.UserDto {
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		log.Debug("Saved user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) GetUserIdentity(identity security.Identity) single.Single[userdtos.UserIdentityDto] {
	userSrc := u.GetByAuthId(identity.GetAuthId())
	return single.Map(userSrc, func(userDto userdtos.UserDto) userdtos.UserIdentityDto {
		userIdentityDto := userdtos.UserIdentityDto{}
		mappers.MapToUserDtoAndIdentityToUserIdentityDto(userDto, identity, &userIdentityDto)
		return userIdentityDto
	})

}

func (u userServiceImpl) GetByAuthId(authId string) single.Single[userdtos.UserDto] {
	userSearchSrc := u.userRepository.FindByAuthId(u.dbHandler.GetCtx(), authId)
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserDto {
		user := userMaybe.OrElse(models.User{})
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		return userDto
	})
}

func NewUserService(
	dbHandler database.DBHandler,
	userRepository repositories.UserRepository,
	userBr businessrules.UserBr,
	errorService apperrors.ErrorService,
) UserService {
	return &userServiceImpl{
		dbHandler:      dbHandler,
		userRepository: userRepository,
		userBr:         userBr,
		errorService:   errorService,
	}
}
