package services

import (
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/businessrules"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/mappers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dbaccess"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errormgmt"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	AddUser(identity security.Identity,
		userSaveDto *userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto)
	UpdateUser(identity security.Identity,
		userSaveDto *userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto)
	GetByAuthId(tokenId string) (*userdtos.UserDto, *errordtos.ErrorResponseDto)
	GetUserIdentity(identity security.Identity) (*userdtos.UserIdentityDto, *errordtos.ErrorResponseDto)
}

type userServiceImpl struct {
	dbClient          dbaccess.DBClient
	transactionRunner dbaccess.TransactionRunner
	userRepository    repositories.UserRepository
	userBr            businessrules.UserBr
}

func (u userServiceImpl) AddUser(identity security.Identity, userSaveDto *userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto) {
	user := &models.User{}
	userDto := &userdtos.UserDto{}

	if errRes := u.userBr.ValidateUserCreate(u.dbClient.GetCtx(), identity, userSaveDto); errRes != nil {
		return nil, errRes
	}
	mappers.MapUserSaveDtoToUser(userSaveDto, user)
	user.AuthId = identity.GetAuthId()

	if err := u.userRepository.Create(u.dbClient.GetCtx(), user); err != nil {
		return nil, errormgmt.CreateInternalErrResponseWithErrLog(err)
	}

	mappers.MapUserToUserDto(user, userDto)
	log.Info("Created user ", userDto)

	return userDto, nil
}

func (u userServiceImpl) UpdateUser(identity security.Identity, userSaveDto *userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto) {
	user := &models.User{}
	userDto := &userdtos.UserDto{}

	if err := u.userRepository.FindByAuthId(
		u.dbClient.GetCtx(), identity.GetAuthId(), user); err != nil {
		return nil, u.dbClient.NotFoundOrElseInternalErrResponse(err)
	}

	if errRes := u.userBr.ValidateUserUpdate(u.dbClient.GetCtx(), userSaveDto, user); errRes != nil {
		return nil, errRes
	}

	mappers.MapUserSaveDtoToUser(userSaveDto, user)

	if err := u.userRepository.Update(u.dbClient.GetCtx(), user); err != nil {
		return nil, errormgmt.CreateInternalErrResponseWithErrLog(err)
	}

	mappers.MapUserToUserDto(user, userDto)
	return userDto, nil
}

func (u userServiceImpl) GetUserIdentity(
	identity security.Identity) (*userdtos.UserIdentityDto, *errordtos.ErrorResponseDto) {
	userDto, errResponse := u.GetByAuthId(identity.GetAuthId())
	if errResponse != nil {
		return nil, errResponse
	}
	userIdentityDto := &userdtos.UserIdentityDto{}
	mappers.MapToUserDtoAndIdentityToUserIdentityDto(userDto, identity, userIdentityDto)
	return userIdentityDto, nil

}

func (u userServiceImpl) GetByAuthId(authId string) (*userdtos.UserDto, *errordtos.ErrorResponseDto) {
	user := &models.User{}
	userDto := &userdtos.UserDto{}

	if err := u.userRepository.FindByAuthId(u.dbClient.GetCtx(), authId, user); err != nil {
		if u.dbClient.IsNotFoundError(err) {
			return userDto, nil
		}
		return nil, u.dbClient.NotFoundOrElseInternalErrResponse(err)
	}

	mappers.MapUserToUserDto(user, userDto)

	return userDto, nil
}

func NewUserService(dbClient dbaccess.DBClient, transactionRunner dbaccess.TransactionRunner,
	userRepository repositories.UserRepository, userBr businessrules.UserBr) UserService {
	return &userServiceImpl{
		dbClient:          dbClient,
		transactionRunner: transactionRunner,
		userRepository:    userRepository,
		userBr:            userBr,
	}
}
