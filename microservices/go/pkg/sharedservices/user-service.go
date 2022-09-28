package sharedservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	sBos "github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedbusinessobjects"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	cModels "github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserService interface {
	RequireUser(ctx context.Context, identity security.Identity) single.Single[sBos.UserBo]
	SaveUser(ctx context.Context, distUserDto userdtos.DistributedUserDto) single.Single[sBos.UserBo]
	DeleteUser(ctx context.Context, distUserDto userdtos.DistributedUserDto) single.Single[sBos.UserBo]
}

type userServiceImpl struct {
	userRepository sharedrepos.UserRepository
	extUserService externalservices.ExtUserService
	errorService   ErrorService
}

func (u userServiceImpl) SaveUser(
	ctx context.Context,
	distUserDto userdtos.DistributedUserDto,
) single.Single[sBos.UserBo] {
	logger.Log.Debugf("saving user %v", distUserDto)
	userSavedSrc := u.saveUserDataAndGetModel(ctx, distUserDto.AuthId, distUserDto.User)
	return single.Map(userSavedSrc, func(u cModels.User) sBos.UserBo {
		userBo := sBos.UserBo{}
		sharedmappers.MapUserToUserBo(u, &userBo)
		return userBo
	})
}

func (u userServiceImpl) saveUserDataAndGetModel(
	ctx context.Context,
	authId string,
	userDto userdtos.UserDto,
) single.Single[cModels.User] {
	userFindSrc := u.userRepository.FindByUserId(ctx, userDto.Id)
	return single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[cModels.User]) single.Single[cModels.User] {
			return option.Map(userMaybe, func(user cModels.User) single.Single[cModels.User] {
				sharedmappers.MapAuthIdAndUserDtoToUser(authId, userDto, &user)
				return u.userRepository.Update(ctx, user)
			}).OrElseGet(func() single.Single[cModels.User] {
				user := cModels.User{}
				sharedmappers.MapAuthIdAndUserDtoToUser(authId, userDto, &user)
				return u.userRepository.Create(ctx, user)
			})
		},
	)
}

func (u userServiceImpl) DeleteUser(
	ctx context.Context,
	distUserDto userdtos.DistributedUserDto,
) single.Single[sBos.UserBo] {
	logger.Log.Debugf("deleting user %v", distUserDto)
	userFindSrc := u.userRepository.FindByUserId(ctx, distUserDto.User.Id)
	userSavedSrc := single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[cModels.User]) single.Single[cModels.User] {
			return option.Map(userMaybe, func(user cModels.User) single.Single[cModels.User] {
				return u.userRepository.Delete(ctx, user)
			}).OrElseGet(func() single.Single[cModels.User] {
				user := cModels.User{}
				sharedmappers.MapDistUserToUser(distUserDto, &user)
				return single.Just(user)
			})
		},
	)
	return single.Map(userSavedSrc, func(u cModels.User) sBos.UserBo {
		userBo := sBos.UserBo{}
		sharedmappers.MapUserToUserBo(u, &userBo)
		return userBo
	})
}

func (u userServiceImpl) RequireUser(ctx context.Context, identity security.Identity) single.Single[sBos.UserBo] {
	userFindSrc := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
	userSrc := single.FlatMap(userFindSrc, func(userMaybe option.Maybe[cModels.User]) single.Single[cModels.User] {
		return option.Map(userMaybe, single.Just[cModels.User]).OrElseGet(func() single.Single[cModels.User] {
			// If user is not stored locally in the database
			extUserFindSrc := u.extUserService.GetByAuthId(ctx, identity.GetAuthId())
			return single.FlatMap(extUserFindSrc, func(extUserDto userdtos.UserDto) single.Single[cModels.User] {
				if !extUserDto.Exists {
					userReqFailRuleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeUserRequireFail)
					return single.Error[cModels.User](apperrors.NewBadReqErrorFromRuleError(userReqFailRuleErr))
				}
				return u.saveUserDataAndGetModel(ctx, identity.GetAuthId(), extUserDto)
			})
		})
	})
	return single.Map(userSrc, func(user cModels.User) sBos.UserBo {
		userBo := sBos.UserBo{}
		sharedmappers.MapUserToUserBo(user, &userBo)
		return userBo
	})
}

func NewUserService(
	userRepository sharedrepos.UserRepository,
	extUserService externalservices.ExtUserService,
) UserService {
	return &userServiceImpl{userRepository: userRepository, extUserService: extUserService}
}
