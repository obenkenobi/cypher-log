package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type UserChangeEventService interface {
	HandleUserChangeEventTxn(
		ctx context.Context,
		userEventDto userdtos.UserChangeEventDto,
	) (userdtos.UserChangeEventResponseDto, error)
}

type UserChangeEventServiceImpl struct {
	userService   sharedservices.UserService
	crudDSHandler dshandlers.CrudDSHandler
	noteService   NoteService
}

func (u UserChangeEventServiceImpl) HandleUserChangeEventTxn(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) (userdtos.UserChangeEventResponseDto, error) {
	return dshandlers.Txn(ctx, u.crudDSHandler,
		func(session dshandlers.Session, ctx context.Context) (userdtos.UserChangeEventResponseDto, error) {
			return u.HandleUserChangeEvent(ctx, userEventDto)
		},
	)
}

func (u UserChangeEventServiceImpl) HandleUserChangeEvent(
	ctx context.Context,
	userEventDto userdtos.UserChangeEventDto,
) (userdtos.UserChangeEventResponseDto, error) {
	switch userEventDto.Action {
	case userdtos.UserSave:
		_, err := u.userService.SaveUser(ctx, userEventDto)
		return userdtos.UserChangeEventResponseDto{Discarded: false}, err
	case userdtos.UserDelete:
		_, err := u.userService.DeleteUser(ctx, userEventDto)
		if err != nil {
			return userdtos.UserChangeEventResponseDto{Discarded: false}, err
		}
		_, err = u.noteService.DeleteByUserIdAndGetCount(ctx, userEventDto.Id)
		return userdtos.UserChangeEventResponseDto{Discarded: false}, err
	default:
		return userdtos.UserChangeEventResponseDto{Discarded: true}, nil
	}
}

func NewUserChangeEventServiceImpl(
	userService sharedservices.UserService,
	crudDSHandler dshandlers.CrudDSHandler,
	noteService NoteService,
) *UserChangeEventServiceImpl {
	return &UserChangeEventServiceImpl{
		userService:   userService,
		crudDSHandler: crudDSHandler,
		noteService:   noteService,
	}
}
