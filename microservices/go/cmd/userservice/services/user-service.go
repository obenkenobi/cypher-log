package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/mappers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserService interface {
	AddUser(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) single.Single[userdtos.UserReadDto]
	UpdateUser(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) single.Single[userdtos.UserReadDto]
	DeleteUser(ctx context.Context, identity security.Identity) single.Single[userdtos.UserReadDto]
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto]
	GetById(ctx context.Context, userId string) single.Single[userdtos.UserReadDto]
	GetUserIdentity(ctx context.Context, identity security.Identity) single.Single[userdtos.UserIdentityDto]
}

type userServiceImpl struct {
	userMsgSendService    UserMsgSendService
	crudDSHandler         dshandlers.CrudDSHandler
	userRepository        repositories.UserRepository
	userBr                businessrules.UserBr
	errorService          sharedservices.ErrorService
	authServerMgmtService AuthServerMgmtService
}

func (u userServiceImpl) AddUser(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserReadDto] {
	userCreateValidationSrc := u.userBr.ValidateUserCreate(ctx, identity, userSaveDto)
	userCreateSrc := single.FlatMap(userCreateValidationSrc, func([]apperrors.RuleError) single.Single[models.User] {
		user := models.User{}
		mappers.UserSaveDtoToUser(userSaveDto, &user)
		user.AuthId = identity.GetAuthId()
		return u.userRepository.Create(ctx, user)
	})
	return single.FlatMap(userCreateSrc, func(user models.User) single.Single[userdtos.UserReadDto] {
		userDto := userdtos.UserReadDto{}
		mappers.UserToUserDto(user, &userDto)
		logger.Log.Debugf("Created user %v", userDto)
		return u.sendUserSave(userDto, identity)
	})

}

func (u userServiceImpl) UpdateUser(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserReadDto] {
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
		mappers.UserSaveDtoToUser(userSaveDto, &user)
		return u.userRepository.Update(ctx, user)
	})
	return single.FlatMap(userSavedSrc, func(user models.User) single.Single[userdtos.UserReadDto] {
		userDto := userdtos.UserReadDto{}
		mappers.UserToUserDto(user, &userDto)
		logger.Log.Debug("Saved user ", userDto)

		return u.sendUserSave(userDto, identity)
	})
}

func (u userServiceImpl) DeleteUser(ctx context.Context, identity security.Identity) single.Single[userdtos.UserReadDto] {
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
		return u.userRepository.Delete(ctx, user)
	})
	userDeletedAuthServerSrc := single.FlatMap(
		userDeletedLocalDBSrc,
		func(user models.User) single.Single[models.User] {
			return single.Map(u.authServerMgmtService.DeleteUser(identity.GetAuthId()),
				func(_ bool) models.User { return user })
		},
	)
	return single.FlatMap(userDeletedAuthServerSrc, func(user models.User) single.Single[userdtos.UserReadDto] {
		logger.Log.Debug("Deleted user ", user)
		userDeleteDto := userdtos.DistUserDeleteDto{}
		mappers.UserToDistUserDeleteDto(user, &userDeleteDto)
		return single.Map(
			u.userMsgSendService.UserDeleteSender().Send(userDeleteDto),
			func(_ userdtos.DistUserDeleteDto) userdtos.UserReadDto {
				userDto := userdtos.UserReadDto{}
				mappers.UserToUserDto(user, &userDto)
				return userDto
			},
		)
	})
}

func (u userServiceImpl) GetUserIdentity(
	ctx context.Context,
	identity security.Identity,
) single.Single[userdtos.UserIdentityDto] {
	userSrc := u.GetByAuthId(ctx, identity.GetAuthId())
	return single.Map(userSrc, func(userDto userdtos.UserReadDto) userdtos.UserIdentityDto {
		userIdentityDto := userdtos.UserIdentityDto{}
		sharedmappers.UserReadDtoAndIdentityToUserIdentityDto(userDto, identity, &userIdentityDto)
		return userIdentityDto
	})

}

func (u userServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto] {
	userSearchSrc := u.userRepository.FindByAuthId(ctx, authId)
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserReadDto {
		user := userMaybe.OrElse(models.User{})
		userDto := userdtos.UserReadDto{}
		mappers.UserToUserDto(user, &userDto)
		return userDto
	})
}

func (u userServiceImpl) GetById(ctx context.Context, userId string) single.Single[userdtos.UserReadDto] {
	userSearchSrc := u.userRepository.FindById(ctx, userId)
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserReadDto {
		user := userMaybe.OrElse(models.User{})
		userDto := userdtos.UserReadDto{}
		mappers.UserToUserDto(user, &userDto)
		return userDto
	})
}

func (u userServiceImpl) sendUserSave(
	userDto userdtos.UserReadDto,
	identity security.Identity,
) single.Single[userdtos.UserReadDto] {
	distUserDto := userdtos.DistUserSaveDto{}
	sharedmappers.UserDtoAndIdentityToDistUserSaveDto(userDto, identity, &distUserDto)
	return single.Map(
		u.userMsgSendService.UserSaveSender().Send(distUserDto),
		func(_ userdtos.DistUserSaveDto) userdtos.UserReadDto {
			return userDto
		},
	)
}

func NewUserService(
	userMsgSendService UserMsgSendService,
	crudDBHandler dshandlers.CrudDSHandler,
	userRepository repositories.UserRepository,
	userBr businessrules.UserBr,
	errorService sharedservices.ErrorService,
	authServerMgmtService AuthServerMgmtService,
) UserService {
	return &userServiceImpl{
		userMsgSendService:    userMsgSendService,
		crudDSHandler:         crudDBHandler,
		userRepository:        userRepository,
		userBr:                userBr,
		errorService:          errorService,
		authServerMgmtService: authServerMgmtService,
	}
}
