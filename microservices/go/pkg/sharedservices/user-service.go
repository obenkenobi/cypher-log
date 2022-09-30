package sharedservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	cModels "github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded/embeddeduser"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserService interface {
	RequireUser(ctx context.Context, identity security.Identity) single.Single[userbos.UserBo]
	SaveUser(ctx context.Context, distUserDto userdtos.DistUserSaveDto) single.Single[userbos.UserBo]
	DeleteUser(ctx context.Context, distUserDto userdtos.DistUserDeleteDto) single.Single[userbos.UserBo]
}

type userServiceImpl struct {
	userRepository sharedrepos.UserRepository
	extUserService externalservices.ExtUserService
	errorService   ErrorService
}

func (u userServiceImpl) SaveUser(
	ctx context.Context,
	distUserDto userdtos.DistUserSaveDto,
) single.Single[userbos.UserBo] {
	logger.Log.Debugf("saving user %v", distUserDto)
	userSavedSrc := u.saveUserDataAndGetModel(ctx, distUserDto.AuthId, distUserDto.BaseUserPublicDto)
	return single.Map(userSavedSrc, func(u cModels.User) userbos.UserBo {
		userBo := userbos.UserBo{}
		sharedmappers.UserModelToUserBo(u, &userBo)
		return userBo
	})
}

func (u userServiceImpl) DeleteUser(
	ctx context.Context,
	distUserDto userdtos.DistUserDeleteDto,
) single.Single[userbos.UserBo] {
	logger.Log.Debugf("deleting user %v", distUserDto)
	userFindSrc := u.userRepository.FindByUserId(ctx, distUserDto.Id)
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

func (u userServiceImpl) RequireUser(ctx context.Context, identity security.Identity) single.Single[userbos.UserBo] {
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

func (u userServiceImpl) saveUserDataAndGetModel(
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

func NewUserService(
	userRepository sharedrepos.UserRepository,
	extUserService externalservices.ExtUserService,
	errorService ErrorService,
) UserService {
	return &userServiceImpl{userRepository: userRepository, extUserService: extUserService, errorService: errorService}
}
