package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
)

type UserBr interface {
	ValidateUserCreate(dbctx context.Context, identity security.Identity,
		dto *userdtos.UserSaveDto) *errordtos.ErrorResponseDto
	ValidateUserUpdate(dbctx context.Context, dto *userdtos.UserSaveDto,
		existing *models.User) *errordtos.ErrorResponseDto
}

type UserBrImpl struct {
	dbHandler      database.DBHandler
	userRepository repositories.UserRepository
}

func (u UserBrImpl) ValidateUserCreate(dbctx context.Context, identity security.Identity,
	dto *userdtos.UserSaveDto) *errordtos.ErrorResponseDto {
	var errCodes []string

	authId := identity.GetAuthId()
	if err := u.userRepository.FindByAuthId(dbctx, authId, &models.User{}); err == nil {
		errCodes = append(errCodes, apperrors.ErrCodeUserAlreadyCreated)
	} else if !u.dbHandler.IsNotFoundError(err) {
		return apperrors.CreateInternalErrResponse(err)
	}

	if returnCodes, err := u.validateUserNameNotTaken(dbctx, dto); err != nil {
		return apperrors.CreateInternalErrResponse(err)
	} else {
		errCodes = append(errCodes, returnCodes...)
	}

	if len(errCodes) == 0 {
		return nil
	}
	return apperrors.CreateErrorResponseFromErrorCodes(errCodes...)
}

func (u UserBrImpl) ValidateUserUpdate(dbctx context.Context, dto *userdtos.UserSaveDto,
	existing *models.User) *errordtos.ErrorResponseDto {
	var errCodes []string

	if dto.UserName != existing.UserName { // If username has changed, ensure the user is not taken
		if returnCodes, err := u.validateUserNameNotTaken(dbctx, dto); err != nil {
			return apperrors.CreateInternalErrResponse(err)
		} else {
			errCodes = append(errCodes, returnCodes...)
		}
	}

	if len(errCodes) == 0 {
		return nil
	}
	return apperrors.CreateErrorResponseFromErrorCodes(errCodes...)
}

func (u UserBrImpl) validateUserNameNotTaken(dbctx context.Context, dto *userdtos.UserSaveDto) ([]string, error) {
	var errCodes []string
	if err := u.userRepository.FindByUsername(dbctx, dto.UserName, &models.User{}); err == nil {
		errCodes = append(errCodes, apperrors.ErrCodeUsernameTaken)
	} else if !u.dbHandler.IsNotFoundError(err) {
		return errCodes, err
	}
	return errCodes, nil
}

func NewUserBrImpl(dbHandler database.DBHandler, userRepository repositories.UserRepository) *UserBrImpl {
	return &UserBrImpl{dbHandler: dbHandler, userRepository: userRepository}
}
