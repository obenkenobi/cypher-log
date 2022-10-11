package services

import (
	"context"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type UserChangeEventService interface {
	HandleUserChangeEventTransaction(
		ctx context.Context,
		userEventDto userdtos.UserChangeEventDto,
	) single.Single[userdtos.UserChangeEventResponseDto]
}

type UserChangeEventServiceImpl struct {
	userService                sharedservices.UserService
	userKeyGeneratorRepository repositories.UserKeyGeneratorRepository
	crudDSHandler              dshandlers.CrudDSHandler
}

func (u UserChangeEventServiceImpl) HandleUserChangeEventTransaction(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) single.Single[userdtos.UserChangeEventResponseDto] {
	return dshandlers.TransactionalSingle(
		ctx,
		u.crudDSHandler,
		func(session dshandlers.Session, ctx context.Context) single.Single[userdtos.UserChangeEventResponseDto] {
			var userResSrc single.Single[userdtos.UserChangeEventResponseDto]
			switch userEventDto.Action {
			case userdtos.UserSave:
				saveUserSrc := u.userService.SaveUser(ctx, userEventDto)
				userResSrc = single.Map(saveUserSrc, func(a userbos.UserBo) userdtos.UserChangeEventResponseDto {
					return userdtos.UserChangeEventResponseDto{Discarded: false}
				})
			case userdtos.UserDelete:
				userDeleteSrc := u.userService.DeleteUser(ctx, userEventDto)
				userKeyDeleteSrc := u.userKeyGeneratorRepository.DeleteByUserIdAndGetCount(ctx, userEventDto.Id)
				userResSrc = single.Map(
					single.Zip2(userDeleteSrc, userKeyDeleteSrc),
					func(_ tuple.T2[userbos.UserBo, int64]) userdtos.UserChangeEventResponseDto {
						return userdtos.UserChangeEventResponseDto{Discarded: false}
					},
				)
			default:
				userResSrc = single.Just(userdtos.UserChangeEventResponseDto{Discarded: true})
			}
			return userResSrc
		},
	)

}

func NewUserChangeEventServiceImpl(
	userService sharedservices.UserService,
	userKeyRepository repositories.UserKeyGeneratorRepository,
	crudDSHandler dshandlers.CrudDSHandler,
) *UserChangeEventServiceImpl {
	return &UserChangeEventServiceImpl{
		userService:                userService,
		userKeyGeneratorRepository: userKeyRepository,
		crudDSHandler:              crudDSHandler,
	}
}