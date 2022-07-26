package sharedservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded/embeddeduser"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	cModels "github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserService interface {
	RequireUser(ctx context.Context, identity security.Identity) single.Single[userbos.UserBo]
	SaveUser(ctx context.Context, userEventDto userdtos.UserChangeEventDto) single.Single[userbos.UserBo]
	DeleteUser(ctx context.Context, userEventDto userdtos.UserChangeEventDto) single.Single[userbos.UserBo]
	UserExistsWithId(ctx context.Context, userId string) single.Single[bool]
}

type UserServiceImpl struct {
	userRepository sharedrepos.UserRepository
	extUserService externalservices.ExtUserService
	errorService   ErrorService
}

func (u UserServiceImpl) SaveUser(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) single.Single[userbos.UserBo] {
	logger.Log.Debugf("saving user %v", userEventDto)
	userSavedSrc := u.saveUserDataAndGetModel(ctx, userEventDto.AuthId, userEventDto.BaseUserPublicDto)
	return single.Map(userSavedSrc, func(u cModels.User) userbos.UserBo {
		userBo := userbos.UserBo{}
		sharedmappers.UserModelToUserBo(u, &userBo)
		return userBo
	})
}

func (u UserServiceImpl) DeleteUser(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) single.Single[userbos.UserBo] {
	logger.Log.Debugf("deleting user %v", userEventDto)
	userFindSrc := u.userRepository.FindByUserId(ctx, userEventDto.Id)
	userSavedSrc := single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[cModels.User]) single.Single[cModels.User] {
			if user, isPresent := userMaybe.Get(); isPresent {
				return u.userRepository.Delete(ctx, user)
			} else {
				return single.Just(cModels.User{})
			}
		},
	)
	return single.Map(userSavedSrc, func(u cModels.User) userbos.UserBo {
		userBo := userbos.UserBo{}
		sharedmappers.UserModelToUserBo(u, &userBo)
		return userBo
	})
}

func (u UserServiceImpl) RequireUser(ctx context.Context, identity security.Identity) single.Single[userbos.UserBo] {
	userFindSrc := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
	userSrc := single.FlatMap(userFindSrc, func(userMaybe option.Maybe[cModels.User]) single.Single[cModels.User] {
		if user, isPresent := userMaybe.Get(); isPresent {
			return single.Just(user)
		} else {
			// If user is not stored locally in the database
			extUserFindSrc := u.extUserService.GetByAuthId(ctx, identity.GetAuthId())
			return single.FlatMap(extUserFindSrc, func(extUserDto userdtos.UserReadDto) single.Single[cModels.User] {
				if !extUserDto.Exists {
					userReqFailRuleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeUserRequireFail)
					return single.Error[cModels.User](apperrors.NewBadReqErrorFromRuleError(userReqFailRuleErr))
				}
				return u.saveUserDataAndGetModel(ctx, identity.GetAuthId(), extUserDto.BaseUserPublicDto)
			})
		}
	})
	return single.Map(userSrc, func(user cModels.User) userbos.UserBo {
		userBo := userbos.UserBo{}
		sharedmappers.UserModelToUserBo(user, &userBo)
		return userBo
	})
}

func (u UserServiceImpl) UserExistsWithId(ctx context.Context, userId string) single.Single[bool] {
	userFindSrc := u.userRepository.FindByUserId(ctx, userId)
	return single.FlatMap(userFindSrc, func(userMaybe option.Maybe[cModels.User]) single.Single[bool] {
		return option.Map(userMaybe, func(_ cModels.User) single.Single[bool] {
			return single.Just(true)
		}).OrElseGet(func() single.Single[bool] {
			extUserFindSrc := u.extUserService.GetById(ctx, userId)
			return single.FlatMap(extUserFindSrc, func(extUserDto userdtos.UserReadDto) single.Single[bool] {
				return single.Just(extUserDto.Exists)
			})
		})
	})
}

func (u UserServiceImpl) saveUserDataAndGetModel(
	ctx context.Context,
	authId string,
	userPublicDto embeddeduser.BaseUserPublicDto,
) single.Single[cModels.User] {
	userFindSrc := u.userRepository.FindByUserId(ctx, userPublicDto.Id)
	return single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[cModels.User]) single.Single[cModels.User] {
			user, isPresent := userMaybe.Get()
			if !isPresent {
				user = cModels.User{}
			}
			sharedmappers.AuthIdAndUserPublicDtoToUserModel(authId, userPublicDto, &user)
			if isPresent {
				return u.userRepository.Update(ctx, user)
			} else {
				return u.userRepository.Create(ctx, user)
			}
		},
	)
}

func NewUserServiceImpl(
	userRepository sharedrepos.UserRepository,
	extUserService externalservices.ExtUserService,
	errorService ErrorService,
) *UserServiceImpl {
	return &UserServiceImpl{userRepository: userRepository, extUserService: extUserService, errorService: errorService}
}
