package services

import (
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/mappers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	AddUser(identity security.Identity, userSaveDto *userdtos.UserSaveDto) stream.Observable[*userdtos.UserDto]
	UpdateUser(identity security.Identity, userSaveDto *userdtos.UserSaveDto) stream.Observable[*userdtos.UserDto]
	GetByAuthId(tokenId string) stream.Observable[*userdtos.UserDto]
	GetUserIdentity(identity security.Identity) stream.Observable[*userdtos.UserIdentityDto]
}

type userServiceImpl struct {
	dbHandler      database.DBHandler
	userRepository repositories.UserRepository
	userBr         businessrules.UserBr
	errorService   errors.ErrorService
}

func (u userServiceImpl) AddUser(
	identity security.Identity,
	userSaveDto *userdtos.UserSaveDto,
) stream.Observable[*userdtos.UserDto] {
	userCreateValidationX := u.userBr.ValidateUserCreate(u.dbHandler.GetCtx(), identity, userSaveDto)
	userCreateX := stream.FlatMap(userCreateValidationX, func([]errors.RuleError) stream.Observable[*models.User] {
		user := &models.User{}
		mappers.MapUserSaveDtoToUser(userSaveDto, user)
		user.AuthId = identity.GetAuthId()
		return u.userRepository.Create(u.dbHandler.GetCtx(), user)
	})
	return stream.Map(userCreateX, func(user *models.User) *userdtos.UserDto {
		userDto := &userdtos.UserDto{}
		mappers.MapUserToUserDto(user, userDto)
		log.Info("Created user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) UpdateUser(
	identity security.Identity,
	userSaveDto *userdtos.UserSaveDto,
) stream.Observable[*userdtos.UserDto] {
	userSearchX := u.userRepository.FindByAuthId(u.dbHandler.GetCtx(), identity.GetAuthId())
	userExistsX := stream.FlatMap(
		userSearchX,
		func(userMaybe option.Maybe[*models.User]) stream.Observable[*models.User] {
			return option.Map(userMaybe, func(user *models.User) stream.Observable[*models.User] {
				return stream.Just(user)
			}).OrElseGet(func() stream.Observable[*models.User] {
				err := errors.NewBadReqErrorFromRuleError(
					u.errorService.RuleErrorFromCode(errors.ErrCodeReqItemsNotFound))
				return stream.Error[*models.User](err)
			})
		},
	)
	userValidatedX := stream.FlatMap(userExistsX, func(existingUser *models.User) stream.Observable[*models.User] {
		validationX := u.userBr.ValidateUserUpdate(u.dbHandler.GetCtx(), userSaveDto, existingUser)
		return stream.Map(validationX, func([]errors.RuleError) *models.User {
			return existingUser
		})
	})
	userSavedX := stream.FlatMap(userValidatedX, func(user *models.User) stream.Observable[*models.User] {
		mappers.MapUserSaveDtoToUser(userSaveDto, user)
		return u.userRepository.Update(u.dbHandler.GetCtx(), user)
	})
	return stream.Map(userSavedX, func(user *models.User) *userdtos.UserDto {
		userDto := &userdtos.UserDto{}
		mappers.MapUserToUserDto(user, userDto)
		log.Info("Created user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) GetUserIdentity(identity security.Identity) stream.Observable[*userdtos.UserIdentityDto] {
	userX := u.GetByAuthId(identity.GetAuthId())
	return stream.Map(userX, func(userDto *userdtos.UserDto) *userdtos.UserIdentityDto {
		userIdentityDto := &userdtos.UserIdentityDto{}
		mappers.MapToUserDtoAndIdentityToUserIdentityDto(userDto, identity, userIdentityDto)
		return userIdentityDto
	})

}

func (u userServiceImpl) GetByAuthId(authId string) stream.Observable[*userdtos.UserDto] {
	userSearchX := u.userRepository.FindByAuthId(u.dbHandler.GetCtx(), authId)
	return stream.Map(userSearchX, func(userMaybe option.Maybe[*models.User]) *userdtos.UserDto {
		user := userMaybe.OrElse(&models.User{})
		userDto := &userdtos.UserDto{}
		mappers.MapUserToUserDto(user, userDto)
		return userDto
	})
}

func NewUserService(
	dbHandler database.DBHandler,
	userRepository repositories.UserRepository,
	userBr businessrules.UserBr,
	errorService errors.ErrorService,
) UserService {
	return &userServiceImpl{
		dbHandler:      dbHandler,
		userRepository: userRepository,
		userBr:         userBr,
		errorService:   errorService,
	}
}
