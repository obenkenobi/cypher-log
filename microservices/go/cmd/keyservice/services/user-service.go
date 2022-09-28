package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/bos"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/mappers"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserService interface {
	RequireUser(ctx context.Context, authId string) single.Single[bos.UserBo]
	SaveUser(ctx context.Context, distUserDto userdtos.DistributedUserDto) single.Single[bos.UserBo]
	DeleteUser(ctx context.Context, distUserDto userdtos.DistributedUserDto) single.Single[bos.UserBo]
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
	extUserService externalservices.ExtUserService
	errorService   errorservices.ErrorService
}

func (u userServiceImpl) SaveUser(
	ctx context.Context,
	distUserDto userdtos.DistributedUserDto,
) single.Single[bos.UserBo] {
	logger.Log.Debugf("saving user %v", distUserDto)
	userSavedSrc := u.saveUserDataAndGetModel(ctx, distUserDto.AuthId, distUserDto.User)
	return single.Map(userSavedSrc, func(u models.User) bos.UserBo {
		userBo := bos.UserBo{}
		mappers.MapUserToUserBo(u, &userBo)
		return userBo
	})
}

func (u userServiceImpl) saveUserDataAndGetModel(
	ctx context.Context,
	authId string,
	userDto userdtos.UserDto,
) single.Single[models.User] {
	userFindSrc := u.userRepository.FindByUserId(ctx, userDto.Id)
	return single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[models.User]) single.Single[models.User] {
			return option.Map(userMaybe, func(user models.User) single.Single[models.User] {
				mappers.MapAuthIdAndUserDtoToUser(authId, userDto, &user)
				return u.userRepository.Update(ctx, user)
			}).OrElseGet(func() single.Single[models.User] {
				user := models.User{}
				mappers.MapAuthIdAndUserDtoToUser(authId, userDto, &user)
				return u.userRepository.Create(ctx, user)
			})
		},
	)
}

func (u userServiceImpl) DeleteUser(
	ctx context.Context,
	distUserDto userdtos.DistributedUserDto,
) single.Single[bos.UserBo] {
	logger.Log.Debugf("deleting user %v", distUserDto)
	userFindSrc := u.userRepository.FindByUserId(ctx, distUserDto.User.Id)
	userSavedSrc := single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[models.User]) single.Single[models.User] {
			return option.Map(userMaybe, func(user models.User) single.Single[models.User] {
				return u.userRepository.Delete(ctx, user)
			}).OrElseGet(func() single.Single[models.User] {
				user := models.User{}
				mappers.MapDistUserToUser(distUserDto, &user)
				return single.Just(user)
			})
		},
	)
	return single.Map(userSavedSrc, func(u models.User) bos.UserBo {
		userBo := bos.UserBo{}
		mappers.MapUserToUserBo(u, &userBo)
		return userBo
	})
}

func (u userServiceImpl) RequireUser(ctx context.Context, authId string) single.Single[bos.UserBo] {
	userFindSrc := u.userRepository.FindByAuthId(ctx, authId)
	userSrc := single.FlatMap(userFindSrc, func(userMaybe option.Maybe[models.User]) single.Single[models.User] {
		return option.Map(userMaybe, single.Just[models.User]).OrElseGet(func() single.Single[models.User] {
			// If user is not stored locally in the database
			extUserFindSrc := u.extUserService.GetByAuthId(ctx, authId)
			return single.FlatMap(extUserFindSrc, func(extUserDto userdtos.UserDto) single.Single[models.User] {
				if !extUserDto.Exists {
					userReqFailRuleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeUserRequireFail)
					return single.Error[models.User](apperrors.NewBadReqErrorFromRuleError(userReqFailRuleErr))
				}
				return u.saveUserDataAndGetModel(ctx, authId, extUserDto)
			})
		})
	})
	return single.Map(userSrc, func(user models.User) bos.UserBo {
		userBo := bos.UserBo{}
		mappers.MapUserToUserBo(user, &userBo)
		return userBo
	})
}

func NewUserService(
	userRepository repositories.UserRepository,
	extUserService externalservices.ExtUserService,
) UserService {
	return &userServiceImpl{userRepository: userRepository, extUserService: extUserService}
}
