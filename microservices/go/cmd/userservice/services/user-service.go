package services

import (
	"context"
	"github.com/joamaki/goreactive/stream"
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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
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
	StartDeleteUser(ctx context.Context, identity security.Identity) single.Single[userdtos.UserReadDto]
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto]
	GetById(ctx context.Context, userId string) single.Single[userdtos.UserReadDto]
	GetUserIdentity(ctx context.Context, identity security.Identity) single.Single[userdtos.UserIdentityDto]
	UsersChangeTask(ctx context.Context)
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
	_ context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserReadDto] {
	return dshandlers.TransactionalSingle(
		u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserReadDto] {
			userCreateValidationSrc := u.userBr.ValidateUserCreate(ctx, identity, userSaveDto)
			userCreateSrc := single.FlatMap(
				userCreateValidationSrc,
				func([]apperrors.RuleError) single.Single[models.User] {
					user := models.User{}
					mappers.UserSaveDtoToUser(userSaveDto, &user)
					user.AuthId = identity.GetAuthId()
					user.Distributed = false
					user.ToBeDeleted = false
					return u.userRepository.Create(ctx, user)
				},
			)
			return single.Map(userCreateSrc, func(user models.User) userdtos.UserReadDto {
				logger.Log.Debug("Saved user ", user)
				return userToUserReadDto(user)
			})
		},
	)

}

func (u userServiceImpl) UpdateUser(
	_ context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) single.Single[userdtos.UserReadDto] {
	return dshandlers.TransactionalSingle(
		u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserReadDto] {
			userSearchSrc := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
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
			userValidatedSrc := single.FlatMap(
				userExistsSrc,
				func(existingUser models.User) single.Single[models.User] {
					validationSrc := u.userBr.ValidateUserUpdate(ctx, userSaveDto, existingUser)
					return single.Map(validationSrc, func([]apperrors.RuleError) models.User { return existingUser })
				},
			)
			userSavedSrc := single.FlatMap(userValidatedSrc, func(user models.User) single.Single[models.User] {
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

func (u userServiceImpl) StartDeleteUser(_ context.Context, identity security.Identity) single.Single[userdtos.UserReadDto] {
	return dshandlers.TransactionalSingle(
		u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserReadDto] {
			userSearchSrc := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
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
			userToBeDeletedSrc := single.FlatMap(userExistsSrc, func(user models.User) single.Single[models.User] {
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
	userSearchSrc := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, authId)
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserReadDto {
		user := userMaybe.OrElse(models.User{})
		return userToUserReadDto(user)
	})
}

func (u userServiceImpl) GetById(ctx context.Context, userId string) single.Single[userdtos.UserReadDto] {
	userSearchSrc := u.userRepository.FindById(ctx, userId)
	return single.Map(userSearchSrc, func(userMaybe option.Maybe[models.User]) userdtos.UserReadDto {
		user := userMaybe.Filter(models.User.WillNotDeleted).OrElse(models.User{})
		return userToUserReadDto(user)
	})
}

func (u userServiceImpl) UsersChangeTask(ctx context.Context) {
	userSampleSrc := u.userRepository.SampleUndistributedUsers(ctx, 100)
	usersCh, errCh := stream.ToChannels(ctx, userSampleSrc)
	var actionSingles []single.Single[any]
	for user := range usersCh {
		var src single.Single[any]
		if user.ToBeDeleted {
			src = single.Map(u.deleteUserTransaction(user), utils.CastToAny[userdtos.UserChangeEventDto])
		} else {
			src = single.Map(u.distributeUserChange(user), utils.CastToAny[userdtos.UserChangeEventDto])
		}
		actionSingles = append(actionSingles, src.ScheduleEagerAsync(ctx))
	}
	err := <-errCh
	if err != nil {
		logger.Log.Error(err)
	}
	for _, actionSrc := range actionSingles {
		if _, err := single.RetrieveValue(ctx, actionSrc); err != nil {
			logger.Log.Error(err)
		}
	}

}

func (u userServiceImpl) deleteUserTransaction(user models.User) single.Single[userdtos.UserChangeEventDto] {
	return dshandlers.TransactionalSingle(
		u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) single.Single[userdtos.UserChangeEventDto] {
			sendUserChangeSrc := u.sendUserChange(user, userdtos.UserDelete).ScheduleEagerAsync(ctx)
			userDeletedLocalDBSrc := u.userRepository.Delete(ctx, user).ScheduleEagerAsync(ctx)
			userDeletedInAppSrc := single.Map(
				single.Zip2(userDeletedLocalDBSrc, sendUserChangeSrc),
				func(t stream.Tuple2[models.User, userdtos.UserChangeEventDto]) models.User {
					return t.V1
				})

			userDeletedAuthServerSrc := single.FlatMap(
				userDeletedInAppSrc,
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

func (u userServiceImpl) distributeUserChange(user models.User) single.Single[userdtos.UserChangeEventDto] {
	return dshandlers.TransactionalSingle(
		u.crudDSHandler,
		func(session dshandlers.Session, ctx context.Context) single.Single[userdtos.UserChangeEventDto] {
			sendUserChangeSrc := u.sendUserChange(user, userdtos.UserSave)
			return single.FlatMap(
				sendUserChangeSrc,
				func(uce userdtos.UserChangeEventDto) single.Single[userdtos.UserChangeEventDto] {
					user := user
					user.ToBeDeleted = false
					user.Distributed = true
					updateSrc := u.userRepository.Update(ctx, user)
					return single.Map(updateSrc, func(a models.User) userdtos.UserChangeEventDto {
						logger.Log.Debugf("Sent user save event for user %v", a)
						return uce
					})
				},
			)

		},
	)
}

func (u userServiceImpl) sendUserChange(
	user models.User,
	action userdtos.UserChangeAction,
) single.Single[userdtos.UserChangeEventDto] {
	distUserDto := userdtos.UserChangeEventDto{}
	mappers.UserToUserChangeEventDto(user, &distUserDto)
	distUserDto.Action = action
	return u.userMsgSendService.UserSaveSender().Send(distUserDto)
}

func userToUserReadDto(user models.User) userdtos.UserReadDto {
	userDto := userdtos.UserReadDto{}
	mappers.UserToUserDto(user, &userDto)
	return userDto
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
