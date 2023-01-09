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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserService interface {
	AddUserTransaction(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) single.Single[userdtos.UserReadDto]
	UpdateUserTransaction(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) single.Single[userdtos.UserReadDto]
	BeginDeletingUserTransaction(ctx context.Context, identity security.Identity) single.Single[userdtos.UserReadDto]
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto]
	GetById(ctx context.Context, userId string) single.Single[userdtos.UserReadDto]
	GetUserIdentity(ctx context.Context, identity security.Identity) single.Single[userdtos.UserIdentityDto]
	UsersChangeTask(ctx context.Context)
}

type UserServiceImpl struct {
	userMsgSendService    UserMsgSendService
	crudDSHandler         dshandlers.CrudDSHandler
	userRepository        repositories.UserRepository
	userBr                businessrules.UserBr
	errorService          sharedservices.ErrorService
	authServerMgmtService AuthServerMgmtService
}

func (u UserServiceImpl) AddUserTransaction(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserReadDto] {
	return dshandlers.TransactionalSingle(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserReadDto] {
			err := u.userBr.ValidateUserCreate(ctx, identity, userSaveDto)
			if err != nil {
				return single.Error[userdtos.UserReadDto](err)
			}
			userCreateSrc := single.FromSupplierCached(func() (models.User, error) {
				user := models.User{}
				mappers.UserSaveDtoToUser(userSaveDto, &user)
				user.AuthId = identity.GetAuthId()
				user.Distributed = false
				user.ToBeDeleted = false
				return u.userRepository.Create(ctx, user)
			})
			return single.Map(userCreateSrc, func(user models.User) userdtos.UserReadDto {
				logger.Log.Debug("Saved user ", user)
				return userToUserReadDto(user)
			})
		},
	)
}

func (u UserServiceImpl) UpdateUserTransaction(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserReadDto] {
	return dshandlers.TransactionalSingle(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserReadDto] {
			userSearchSrc := single.FromSupplierCached(func() (option.Maybe[models.User], error) {
				return u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
			})
			userExistsSrc := single.MapWithError(userSearchSrc,
				func(userMaybe option.Maybe[models.User]) (models.User, error) {
					if user, isPresent := userMaybe.Get(); isPresent {
						return user, nil
					} else {
						err := apperrors.NewBadReqErrorFromRuleError(
							u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
						return user, err
					}
				},
			)
			userValidatedSrc := single.MapWithError(userExistsSrc,
				func(existingUser models.User) (models.User, error) {
					err := u.userBr.ValidateUserUpdate(ctx, userSaveDto, existingUser)
					return existingUser, err
				},
			)
			userSavedSrc := single.MapWithError(userValidatedSrc, func(user models.User) (models.User, error) {
				mappers.UserSaveDtoToUser(userSaveDto, &user)
				user.Distributed = false
				user.ToBeDeleted = false
				return u.userRepository.Update(ctx, user)
			})
			return single.Map(userSavedSrc, func(user models.User) userdtos.UserReadDto {
				logger.Log.Debug("Saved user ", user)
				return userToUserReadDto(user)
			})
		},
	)

}

func (u UserServiceImpl) BeginDeletingUserTransaction(
	ctx context.Context,
	identity security.Identity,
) single.Single[userdtos.UserReadDto] {
	return dshandlers.TransactionalSingle(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserReadDto] {
			userSearchSrc := single.FromSupplierCached(func() (option.Maybe[models.User], error) {
				return u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
			})
			userExistsSrc := single.MapWithError(
				userSearchSrc,
				func(userMaybe option.Maybe[models.User]) (models.User, error) {
					if user, isPresent := userMaybe.Get(); isPresent {
						return user, nil
					} else {
						err := apperrors.NewBadReqErrorFromRuleError(
							u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
						return user, err
					}
				},
			)
			userToBeDeletedSrc := single.MapWithError(userExistsSrc, func(user models.User) (models.User, error) {
				user.Distributed = false
				user.ToBeDeleted = true
				return u.userRepository.Update(ctx, user)
			})
			return single.Map(userToBeDeletedSrc, func(user models.User) userdtos.UserReadDto {
				logger.Log.Debug("Starting to delete user ", user)
				return userToUserReadDto(user)
			})
		})
}

func (u UserServiceImpl) GetUserIdentity(
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

func (u UserServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto] {
	userSearchSrc := single.FromSupplierCached(func() (option.Maybe[models.User], error) {
		return u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, authId)
	})
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserReadDto {
		user := userMaybe.OrElse(models.User{})
		return userToUserReadDto(user)
	})
}

func (u UserServiceImpl) GetById(ctx context.Context, userId string) single.Single[userdtos.UserReadDto] {
	userSearchSrc := single.FromSupplierCached(func() (option.Maybe[models.User], error) {
		return u.userRepository.FindById(ctx, userId)
	})
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserReadDto {
		user := userMaybe.Filter(models.User.WillNotDeleted).OrElse(models.User{})
		return userToUserReadDto(user)
	})
}

func (u UserServiceImpl) UsersChangeTask(ctx context.Context) {
	userSample, err := u.userRepository.SampleUndistributedUsers(ctx, 100)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	for _, user := range userSample {
		var err error
		if user.ToBeDeleted {
			_, err = single.RetrieveValue(ctx, u.deleteUserTransaction(ctx, user))
		} else {
			_, err = single.RetrieveValue(ctx, u.distributeUserChangeTransaction(ctx, user))
		}
		if err != nil {
			logger.Log.Error(err)
		}
	}

}

func (u UserServiceImpl) deleteUserTransaction(
	ctx context.Context,
	user models.User,
) single.Single[userdtos.UserChangeEventDto] {
	return dshandlers.TransactionalSingle(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserChangeEventDto] {
			sendUserChangeSrc := u.sendUserChange(user, userdtos.UserDelete)
			userDeletedLocalDBSrc := single.FromSupplierCached(func() (models.User, error) {
				return u.userRepository.Delete(ctx, user)
			})
			userDeletedAuthServerSrc := single.FlatMap(
				userDeletedLocalDBSrc,
				func(user models.User) single.Single[models.User] {
					return single.Map(u.authServerMgmtService.DeleteUser(user.AuthId),
						func(_ bool) models.User { return user })
				},
			)
			return single.FlatMap(
				userDeletedAuthServerSrc,
				func(user models.User) single.Single[userdtos.UserChangeEventDto] {
					logger.Log.Debug("Deleted user ", user)
					return sendUserChangeSrc
				},
			)
		})
}

func (u UserServiceImpl) distributeUserChangeTransaction(
	ctx context.Context,
	user models.User,
) single.Single[userdtos.UserChangeEventDto] {
	return dshandlers.TransactionalSingle(ctx, u.crudDSHandler,
		func(session dshandlers.Session, ctx context.Context) single.Single[userdtos.UserChangeEventDto] {
			sendUserChangeSrc := u.sendUserChange(user, userdtos.UserSave)
			return single.MapWithError(
				sendUserChangeSrc,
				func(event userdtos.UserChangeEventDto) (userdtos.UserChangeEventDto, error) {
					user := user
					user.ToBeDeleted = false
					user.Distributed = true
					updatedUser, err := u.userRepository.Update(ctx, user)
					if err == nil {
						logger.Log.Debugf("Sent user save event for user %v", updatedUser)
					}
					return event, err
				},
			)

		},
	)
}

func (u UserServiceImpl) sendUserChange(
	user models.User,
	action userdtos.UserChangeAction,
) single.Single[userdtos.UserChangeEventDto] {
	distUserDto := userdtos.UserChangeEventDto{}
	mappers.UserToUserChangeEventDto(user, &distUserDto)
	distUserDto.Action = action
	return single.FromSupplierCached(func() (userdtos.UserChangeEventDto, error) {
		return u.userMsgSendService.UserSaveSender().Send(distUserDto)
	})
}

func userToUserReadDto(user models.User) userdtos.UserReadDto {
	userDto := userdtos.UserReadDto{}
	mappers.UserToUserDto(user, &userDto)
	return userDto
}

func NewUserServiceImpl(
	userMsgSendService UserMsgSendService,
	crudDBHandler dshandlers.CrudDSHandler,
	userRepository repositories.UserRepository,
	userBr businessrules.UserBr,
	errorService sharedservices.ErrorService,
	authServerMgmtService AuthServerMgmtService,
) *UserServiceImpl {
	return &UserServiceImpl{
		userMsgSendService:    userMsgSendService,
		crudDSHandler:         crudDBHandler,
		userRepository:        userRepository,
		userBr:                userBr,
		errorService:          errorService,
		authServerMgmtService: authServerMgmtService,
	}
}
