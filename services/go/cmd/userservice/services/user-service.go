package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/mappers"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dbaccess"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	AddUser(userSaveDto userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto)
	UpdateUser(identityHolder security.IdentityHolder, userSaveDto userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto)
	GetByProviderUserId(tokenId string) (*userdtos.UserDto, *errordtos.ErrorResponseDto)
}

type userServiceImpl struct {
	dbClient          dbaccess.DBClient
	transactionRunner dbaccess.TransactionRunner
	userRepository    repositories.UserRepository
}

func (u userServiceImpl) AddUser(userSaveDto userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto) {
	user := &models.User{}
	savedUser := &models.User{}
	userDto := &userdtos.UserDto{}
	var errorResponse *errordtos.ErrorResponseDto = nil

	if err := u.transactionRunner.ExecTransaction(func(session dbaccess.Session, ctx context.Context) error {
		mappers.MapUserSaveDtoToUser(&userSaveDto, user)

		if err := u.userRepository.Create(ctx, user); err != nil {
			errorResponse = errors.CreateInternalErrResponseWithLog(err)
			return session.AbortTransaction(ctx)
		}

		if err := u.userRepository.FindByProviderUserId(
			u.dbClient.GetContext(), user.ID.String(), savedUser); err != nil {
			errorResponse = errors.CreateInternalErrResponseWithLog(err)
			return session.AbortTransaction(ctx)
		}

		mappers.MapUserToUserDto(savedUser, userDto)

		log.Info("Created user", userDto)

		return session.CommitTransaction(ctx)
	}); err != nil {
		errorResponse = errors.CreateInternalErrResponseWithLog(err)
	}
	return userDto, errorResponse
}

func (u userServiceImpl) UpdateUser(identityHolder security.IdentityHolder, userSaveDto userdtos.UserSaveDto) (*userdtos.UserDto, *errordtos.ErrorResponseDto) {
	user := &models.User{}
	savedUser := &models.User{}
	userDto := &userdtos.UserDto{}
	var errorResponse *errordtos.ErrorResponseDto = nil

	if err := u.transactionRunner.ExecTransaction(func(session dbaccess.Session, ctx context.Context) error {

		if err := u.userRepository.FindByProviderUserId(
			u.dbClient.GetContext(), identityHolder.GetIdFromProvider(), user); err != nil {
			errorResponse = errors.CreateInternalErrResponseWithLog(err)
			return session.AbortTransaction(ctx)
		}

		mappers.MapUserSaveDtoToUser(&userSaveDto, user)

		if err := u.userRepository.Update(ctx, user); err != nil {
			errorResponse = errors.CreateInternalErrResponseWithLog(err)
			return session.AbortTransaction(ctx)
		}

		if err := u.userRepository.FindByProviderUserId(
			u.dbClient.GetContext(), user.ID.String(), savedUser); err != nil {
			errorResponse = errors.CreateInternalErrResponseWithLog(err)
			return session.AbortTransaction(ctx)
		}

		mappers.MapUserToUserDto(savedUser, userDto)

		return session.CommitTransaction(ctx)
	}); err != nil {
		errorResponse = errors.CreateInternalErrResponseWithLog(err)
	}
	return userDto, errorResponse
}

func (u userServiceImpl) GetByProviderUserId(providerUserId string) (*userdtos.UserDto, *errordtos.ErrorResponseDto) {
	user := &models.User{}
	userDto := &userdtos.UserDto{}

	if err := u.userRepository.FindByProviderUserId(u.dbClient.GetContext(), providerUserId, user); err != nil {
		return nil, errors.CreateInternalErrResponseWithLog(err)
	}

	mappers.MapUserToUserDto(user, userDto)

	return userDto, nil
}

func NewUserService(dbClient dbaccess.DBClient, transactionRunner dbaccess.TransactionRunner,
	userRepository repositories.UserRepository) UserService {
	return &userServiceImpl{
		dbClient:          dbClient,
		transactionRunner: transactionRunner,
		userRepository:    userRepository,
	}
}
