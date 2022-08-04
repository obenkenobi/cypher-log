package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/mappers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/database/dbservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	AddUser(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) single.Single[userdtos.UserDto]
	UpdateUser(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) single.Single[userdtos.UserDto]
	DeleteUser(ctx context.Context, identity security.Identity) single.Single[userdtos.UserDto]
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto]
	GetUserIdentity(ctx context.Context, identity security.Identity) single.Single[userdtos.UserIdentityDto]
}

type userServiceImpl struct {
	crudDBHandler         dbservices.CrudDBHandler
	userRepository        repositories.UserRepository
	userBr                businessrules.UserBr
	errorService          errorservices.ErrorService
	authServerMgmtService AuthServerMgmtService
}

func (u userServiceImpl) AddUser(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserDto] {
	userCreateValidationSrc := u.userBr.ValidateUserCreate(ctx, identity, userSaveDto)
	userCreateSrc := single.FlatMap(userCreateValidationSrc, func([]apperrors.RuleError) single.Single[models.User] {
		user := models.User{}
		mappers.MapUserSaveDtoToUser(userSaveDto, &user)
		user.AuthId = identity.GetAuthId()
		return single.MapDerefPtr(u.userRepository.Create(ctx, &user))
	})
	return single.Map(userCreateSrc, func(user models.User) userdtos.UserDto {
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		log.Debug("Created user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) UpdateUser(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserDto] {
	userSearchSrc := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
	userExistsSrc := single.MapWithError(
		userSearchSrc,
		func(userMaybe option.Maybe[models.User]) (models.User, error) {
			if user, ok := userMaybe.Get(); ok {
				return user, nil
			} else {
				err := apperrors.NewBadReqErrorFromRuleError(
					u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
				return user, err
			}
		},
	)
	userValidatedSrc := single.FlatMap(userExistsSrc, func(existingUser models.User) single.Single[models.User] {
		validationSrc := u.userBr.ValidateUserUpdate(ctx, userSaveDto, existingUser)
		return single.Map(validationSrc, func([]apperrors.RuleError) models.User { return existingUser })
	})
	userSavedSrc := single.FlatMap(userValidatedSrc, func(user models.User) single.Single[models.User] {
		mappers.MapUserSaveDtoToUser(userSaveDto, &user)
		return single.MapDerefPtr(u.userRepository.Update(ctx, &user))
	})
	return single.Map(userSavedSrc, func(user models.User) userdtos.UserDto {
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		log.Debug("Saved user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) DeleteUser(ctx context.Context, identity security.Identity) single.Single[userdtos.UserDto] {
	userSearchSrc := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
	userExistsSrc := single.MapWithError(
		userSearchSrc,
		func(userMaybe option.Maybe[models.User]) (models.User, error) {
			if user, ok := userMaybe.Get(); ok {
				return user, nil
			} else {
				err := apperrors.NewBadReqErrorFromRuleError(
					u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
				return user, err
			}
		},
	)
	userDeletedLocalDBSrc := single.FlatMap(userExistsSrc, func(user models.User) single.Single[models.User] {
		return single.MapDerefPtr(u.userRepository.Delete(ctx, &user))
	})
	userDeletedAuthServerSrc := single.FlatMap(
		userDeletedLocalDBSrc,
		func(user models.User) single.Single[models.User] {
			return single.Map(u.authServerMgmtService.DeleteUser(identity.GetAuthId()),
				func(_ bool) models.User { return user })
		},
	)
	return single.Map(userDeletedAuthServerSrc, func(user models.User) userdtos.UserDto {
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		log.Debug("Deleted user ", userDto)
		return userDto
	})
}

func (u userServiceImpl) GetUserIdentity(
	ctx context.Context,
	identity security.Identity,
) single.Single[userdtos.UserIdentityDto] {
	userSrc := u.GetByAuthId(ctx, identity.GetAuthId())
	return single.Map(userSrc, func(userDto userdtos.UserDto) userdtos.UserIdentityDto {
		userIdentityDto := userdtos.UserIdentityDto{}
		mappers.MapToUserDtoAndIdentityToUserIdentityDto(userDto, identity, &userIdentityDto)
		return userIdentityDto
	})

}

func (u userServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	userSearchSrc := u.userRepository.FindByAuthId(ctx, authId)
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserDto {
		user := userMaybe.OrElse(models.User{})
		userDto := userdtos.UserDto{}
		mappers.MapUserToUserDto(user, &userDto)
		return userDto
	})
}

func NewUserService(
	crudDBHandler dbservices.CrudDBHandler,
	userRepository repositories.UserRepository,
	userBr businessrules.UserBr,
	errorService errorservices.ErrorService,
	authServerMgmtService AuthServerMgmtService,
) UserService {
	return &userServiceImpl{
		crudDBHandler:         crudDBHandler,
		userRepository:        userRepository,
		userBr:                userBr,
		errorService:          errorService,
		authServerMgmtService: authServerMgmtService,
	}
}
