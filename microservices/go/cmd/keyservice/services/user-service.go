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
	SaveUserAndGetBO(ctx context.Context, userDto userdtos.UserDto) single.Single[bos.UserBo]
	DeleteUserIfFoundAndGetBO(ctx context.Context, userDto userdtos.UserDto) single.Single[bos.UserBo]
}

type userServiceImpl struct {
	userRepository repositories.UserRepository
	extUserService externalservices.ExtUserService
	errorService   errorservices.ErrorService
}

func (u userServiceImpl) SaveUserAndGetBO(ctx context.Context, userDto userdtos.UserDto) single.Single[bos.UserBo] {
	logger.Log.Debugf("saving user %v", userDto)
	userFindSrc := u.userRepository.FindByUserId(ctx, userDto.Id)
	userSavedSrc := single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[models.User]) single.Single[models.User] {
			return option.Map(userMaybe, func(user models.User) single.Single[models.User] {
				mappers.MapUserDtoToUser(userDto, &user)
				return u.userRepository.Update(ctx, user)
			}).OrElseGet(func() single.Single[models.User] {
				user := models.User{}
				mappers.MapUserDtoToUser(userDto, &user)
				return u.userRepository.Create(ctx, user)
			})
		},
	)
	return single.Map(userSavedSrc, func(u models.User) bos.UserBo {
		userBo := bos.UserBo{}
		mappers.MapUserToUserBo(u, &userBo)
		return userBo
	})
}

func (u userServiceImpl) DeleteUserIfFoundAndGetBO(
	ctx context.Context,
	userDto userdtos.UserDto,
) single.Single[bos.UserBo] {
	logger.Log.Debugf("deleting user %v", userDto)
	userFindSrc := u.userRepository.FindByUserId(ctx, userDto.Id)
	userSavedSrc := single.FlatMap(
		userFindSrc,
		func(userMaybe option.Maybe[models.User]) single.Single[models.User] {
			return option.Map(userMaybe, func(user models.User) single.Single[models.User] {
				return u.userRepository.Delete(ctx, user)
			}).OrElseGet(func() single.Single[models.User] {
				user := models.User{}
				mappers.MapUserDtoToUser(userDto, &user)
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
	extUserFindSrc := u.extUserService.GetByAuthId(ctx, authId)
	return single.FlatMap(extUserFindSrc, func(extUserDto userdtos.UserDto) single.Single[bos.UserBo] {
		if !extUserDto.Exists {
			userReqFailRuleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeUserRequireFail)
			return single.Error[bos.UserBo](apperrors.NewBadReqErrorFromRuleError(userReqFailRuleErr))
		}
		return u.SaveUserAndGetBO(ctx, extUserDto)
	})
}

func NewUserService(
	userRepository repositories.UserRepository,
	extUserService externalservices.ExtUserService,
) UserService {
	return &userServiceImpl{userRepository: userRepository, extUserService: extUserService}
}
