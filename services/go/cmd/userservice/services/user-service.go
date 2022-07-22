package services

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dbaccess"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type UserService interface {
	AddUser(userSaveDto userdtos.UserSaveDto) (userdtos.UserDto, *errordtos.ErrorResponseDto)
	UpdateUser(userSaveDto userdtos.UserSaveDto) (userdtos.UserDto, *errordtos.ErrorResponseDto)
	GetByProviderUserId(tokenId string) userdtos.UserDto
}

type userServiceImpl struct {
	dbClient          dbaccess.DBClient
	transactionRunner dbaccess.TransactionRunner
	userRepository    repositories.UserRepository
}

func (u userServiceImpl) AddUser(userSaveDto userdtos.UserSaveDto) (userdtos.UserDto, *errordtos.ErrorResponseDto) {
	userDto := userdtos.UserDto{
		Id:          fmt.Sprintf("%v", rand.Intn(1000000)),
		UserAdded:   true,
		UserName:    userSaveDto.UserName,
		DisplayName: userSaveDto.DisplayName,
	}
	log.Info("Created user", userDto)
	return userDto, nil
}

func (u userServiceImpl) UpdateUser(userSaveDto userdtos.UserSaveDto) (userdtos.UserDto, *errordtos.ErrorResponseDto) {
	userDto := userdtos.UserDto{
		Id:          fmt.Sprintf("%v", rand.Intn(1000000)),
		UserAdded:   true,
		UserName:    userSaveDto.UserName,
		DisplayName: userSaveDto.DisplayName,
	}
	log.Info("Created user", userDto)
	return userDto, nil
}

func (u userServiceImpl) GetByProviderUserId(tokenId string) userdtos.UserDto {
	if tokenId != "1" {
		return userdtos.UserDto{
			Id:          fmt.Sprintf("%v", rand.Intn(1000000)),
			UserAdded:   true,
			UserName:    "number1",
			DisplayName: "robbie",
		}
	} else {
		return userdtos.UserDto{
			UserAdded: false,
		}
	}
}

func NewUserService(dbClient dbaccess.DBClient, transactionRunner dbaccess.TransactionRunner,
	userRepository repositories.UserRepository) UserService {
	return &userServiceImpl{
		dbClient:          dbClient,
		transactionRunner: transactionRunner,
		userRepository:    userRepository,
	}
}
