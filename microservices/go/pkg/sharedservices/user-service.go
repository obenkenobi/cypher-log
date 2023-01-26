package sharedservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded/embeddeduser"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
	cModels "github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedrepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/externalservices"
)

type UserService interface {
	RequireUser(ctx context.Context, identity security.Identity) (userbos.UserBo, error)
	SaveUser(ctx context.Context, userEventDto userdtos.UserChangeEventDto) (userbos.UserBo, error)
	DeleteUser(ctx context.Context, userEventDto userdtos.UserChangeEventDto) (userbos.UserBo, error)
	UserExistsWithId(ctx context.Context, userId string) (bool, error)
}

type UserServiceImpl struct {
	userRepository sharedrepos.UserRepository
	extUserService externalservices.ExtUserService
	errorService   ErrorService
}

func (u UserServiceImpl) SaveUser(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) (userbos.UserBo, error) {
	logger.Log.WithContext(ctx).Debugf("saving user %v", userEventDto)
	user, err := u.saveUserDataAndGetModel(ctx, userEventDto.AuthId, userEventDto.BaseUserPublicDto)
	if err != nil {
		return userbos.UserBo{}, err
	}
	userBo := userbos.UserBo{}
	sharedmappers.UserModelToUserBo(user, &userBo)
	return userBo, nil
}

func (u UserServiceImpl) DeleteUser(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) (userbos.UserBo, error) {
	logger.Log.WithContext(ctx).Debugf("deleting user %v", userEventDto)
	userMaybe, err := u.userRepository.FindByUserId(ctx, userEventDto.Id)
	if err != nil {
		return userbos.UserBo{}, err
	}

	deletedUser := cModels.User{}
	if user, isPresent := userMaybe.Get(); isPresent {
		deletedUser, err = u.userRepository.Delete(ctx, user)
		if err != nil {
			return userbos.UserBo{}, err
		}
	} else {
		return userbos.UserBo{}, nil
	}

	userBo := userbos.UserBo{}
	sharedmappers.UserModelToUserBo(deletedUser, &userBo)
	return userBo, nil
}

func (u UserServiceImpl) RequireUser(ctx context.Context, identity security.Identity) (userbos.UserBo, error) {
	userMaybe, err := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
	if err != nil {
		return userbos.UserBo{}, err
	}
	user, isPresent := userMaybe.Get()

	if !isPresent {
		extUserDto, err := u.extUserService.GetByAuthId(ctx, identity.GetAuthId())
		if err != nil {
			return userbos.UserBo{}, err
		}
		if !extUserDto.Exists {
			userReqFailRuleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeUserRequireFail)
			return userbos.UserBo{}, apperrors.NewBadReqErrorFromRuleError(userReqFailRuleErr)
		}
		user, err = u.saveUserDataAndGetModel(ctx, identity.GetAuthId(), extUserDto.BaseUserPublicDto)
		if err != nil {
			return userbos.UserBo{}, err
		}
	}

	userBo := userbos.UserBo{}
	sharedmappers.UserModelToUserBo(user, &userBo)
	return userBo, nil
}

func (u UserServiceImpl) UserExistsWithId(ctx context.Context, userId string) (bool, error) {
	userMaybe, err := u.userRepository.FindByUserId(ctx, userId)
	if err != nil {
		return false, err
	}

	_, existsInRepo := userMaybe.Get()
	if existsInRepo {
		return true, nil
	}

	extUserDto, err := u.extUserService.GetById(ctx, userId)
	if err != nil {
		return false, err
	}
	return extUserDto.Exists, nil
}

func (u UserServiceImpl) saveUserDataAndGetModel(
	ctx context.Context,
	authId string,
	userPublicDto embeddeduser.BaseUserPublicDto,
) (cModels.User, error) {
	userMaybe, err := u.userRepository.FindByUserId(ctx, userPublicDto.Id)
	if err != nil {
		return cModels.User{}, err
	}
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
}

func NewUserServiceImpl(
	userRepository sharedrepos.UserRepository,
	extUserService externalservices.ExtUserService,
	errorService ErrorService,
) *UserServiceImpl {
	return &UserServiceImpl{userRepository: userRepository, extUserService: extUserService, errorService: errorService}
}
