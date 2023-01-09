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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type UserService interface {
	AddUserTransaction(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) (userdtos.UserReadDto, error)
	UpdateUserTransaction(
		ctx context.Context,
		identity security.Identity,
		userSaveDto userdtos.UserSaveDto,
	) (userdtos.UserReadDto, error)
	BeginDeletingUserTransaction(ctx context.Context, identity security.Identity) (userdtos.UserReadDto, error)
	GetByAuthId(ctx context.Context, authId string) (userdtos.UserReadDto, error)
	GetById(ctx context.Context, userId string) (userdtos.UserReadDto, error)
	GetUserIdentity(ctx context.Context, identity security.Identity) (userdtos.UserIdentityDto, error)
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
) (userdtos.UserReadDto, error) {
	return dshandlers.Transactional(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) (userdtos.UserReadDto, error) {
			err := u.userBr.ValidateUserCreate(ctx, identity, userSaveDto)
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			user := models.User{}
			mappers.UserSaveDtoToUser(userSaveDto, &user)
			user.AuthId = identity.GetAuthId()
			user.Distributed = false
			user.ToBeDeleted = false

			createdUser, err := u.userRepository.Create(ctx, user)
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			logger.Log.Debug("Saved user ", user)
			return userToUserReadDto(createdUser), nil
		},
	)
}

func (u UserServiceImpl) UpdateUserTransaction(
	ctx context.Context,
	identity security.Identity,
	userSaveDto userdtos.UserSaveDto,
) (userdtos.UserReadDto, error) {
	return dshandlers.Transactional(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) (userdtos.UserReadDto, error) {
			userSearch, err := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			user, isPresent := userSearch.Get()
			if !isPresent {
				err := apperrors.NewBadReqErrorFromRuleError(
					u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
				return userdtos.UserReadDto{}, err
			}

			err = u.userBr.ValidateUserUpdate(ctx, userSaveDto, user)
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			mappers.UserSaveDtoToUser(userSaveDto, &user)
			user.Distributed = false
			user.ToBeDeleted = false

			updatedUser, err := u.userRepository.Update(ctx, user)
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			logger.Log.Debug("Saved user ", updatedUser)
			return userToUserReadDto(updatedUser), nil
		},
	)

}

func (u UserServiceImpl) BeginDeletingUserTransaction(
	ctx context.Context,
	identity security.Identity,
) (userdtos.UserReadDto, error) {
	return dshandlers.Transactional(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) (userdtos.UserReadDto, error) {
			userSearch, err := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			user, isPresent := userSearch.Get()
			if !isPresent {
				err := apperrors.NewBadReqErrorFromRuleError(
					u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
				return userdtos.UserReadDto{}, err
			}

			user.Distributed = false
			user.ToBeDeleted = true

			updatedUser, err := u.userRepository.Update(ctx, user)
			if err != nil {
				return userdtos.UserReadDto{}, err
			}

			logger.Log.Debug("Starting to delete user ", updatedUser)
			return userToUserReadDto(updatedUser), nil
		})
}

func (u UserServiceImpl) GetUserIdentity(
	ctx context.Context,
	identity security.Identity,
) (userdtos.UserIdentityDto, error) {
	userDto, err := u.GetByAuthId(ctx, identity.GetAuthId())
	if err != nil {
		return userdtos.UserIdentityDto{}, err
	}
	userIdentityDto := userdtos.UserIdentityDto{}
	sharedmappers.UserReadDtoAndIdentityToUserIdentityDto(userDto, identity, &userIdentityDto)
	return userIdentityDto, nil
}

func (u UserServiceImpl) GetByAuthId(ctx context.Context, authId string) (userdtos.UserReadDto, error) {
	userSearch, err := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, authId)
	if err != nil {
		return userdtos.UserReadDto{}, err
	}
	user := userSearch.OrElse(models.User{})
	return userToUserReadDto(user), nil
}

func (u UserServiceImpl) GetById(ctx context.Context, userId string) (userdtos.UserReadDto, error) {
	userSearch, err := u.userRepository.FindById(ctx, userId)
	if err != nil {
		return userdtos.UserReadDto{}, err
	}
	user := userSearch.Filter(models.User.WillNotBeDeleted).OrElse(models.User{})
	return userToUserReadDto(user), nil
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
			_, err = u.deleteUserTransaction(ctx, user)
		} else {
			_, err = u.distributeUserChangeTransaction(ctx, user)
		}
		if err != nil {
			logger.Log.Error(err)
		}
	}

}

func (u UserServiceImpl) deleteUserTransaction(
	ctx context.Context,
	user models.User,
) (userdtos.UserChangeEventDto, error) {
	return dshandlers.Transactional(ctx, u.crudDSHandler,
		func(s dshandlers.Session, ctx context.Context) (userdtos.UserChangeEventDto, error) {
			event, err := u.sendUserChange(user, userdtos.UserDelete)
			if err != nil {
				return event, err
			}

			deletedUser, err := u.userRepository.Delete(ctx, user)
			if err != nil {
				return event, err
			}

			_, err = u.authServerMgmtService.DeleteUser(deletedUser.AuthId)
			if err != nil {
				return event, err
			}

			logger.Log.Debug("Deleted user ", deletedUser)
			return event, err
		})
}

func (u UserServiceImpl) distributeUserChangeTransaction(
	ctx context.Context,
	user models.User,
) (userdtos.UserChangeEventDto, error) {
	return dshandlers.Transactional(ctx, u.crudDSHandler,
		func(session dshandlers.Session, ctx context.Context) (userdtos.UserChangeEventDto, error) {
			event, err := u.sendUserChange(user, userdtos.UserDelete)
			if err != nil {
				return event, err
			}

			user.ToBeDeleted = false
			user.Distributed = true
			updatedUser, err := u.userRepository.Update(ctx, user)
			if err == nil {
				logger.Log.Debugf("Sent user save event for user %v", updatedUser)
			}
			return event, err
		},
	)
}

func (u UserServiceImpl) sendUserChange(
	user models.User,
	action userdtos.UserChangeAction,
) (userdtos.UserChangeEventDto, error) {
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
