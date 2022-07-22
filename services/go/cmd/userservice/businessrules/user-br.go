package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dbaccess"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errormgmt"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
)

type UserBr interface {
	ValidateUserCreate(dbctx context.Context, identity security.Identity,
		dto *userdtos.UserSaveDto) *errordtos.ErrorResponseDto
	ValidateUserUpdate(dbctx context.Context, dto *userdtos.UserSaveDto,
		existing *models.User) *errordtos.ErrorResponseDto
}

type UserBrImpl struct {
	dbClient       dbaccess.DBClient
	userRepository repositories.UserRepository
}

func (u UserBrImpl) ValidateUserCreate(dbctx context.Context, identity security.Identity,
	dto *userdtos.UserSaveDto) *errordtos.ErrorResponseDto {
	var errCodes []string

	authId := identity.GetAuthId()
	if err := u.userRepository.FindByAuthId(dbctx, authId, &models.User{}); err == nil {
		errCodes = append(errCodes, errormgmt.ErrCodeUserAlreadyCreated)
	} else if !u.dbClient.IsNotFoundError(err) {
		return errormgmt.CreateInternalErrResponseWithErrLog(err)
	}

	if returnCodes, err := u.validateUserNameNotTaken(dbctx, dto); err != nil {
		return errormgmt.CreateInternalErrResponseWithErrLog(err)
	} else {
		errCodes = append(errCodes, returnCodes...)
	}

	if len(errCodes) == 0 {
		return nil
	}
	return errormgmt.CreateErrorResponseFromErrorCodes(errCodes...)
}

func (u UserBrImpl) ValidateUserUpdate(dbctx context.Context, dto *userdtos.UserSaveDto,
	existing *models.User) *errordtos.ErrorResponseDto {
	var errCodes []string

	if dto.UserName != existing.UserName { // If username has changed, ensure the user is not taken
		if returnCodes, err := u.validateUserNameNotTaken(dbctx, dto); err != nil {
			return errormgmt.CreateInternalErrResponseWithErrLog(err)
		} else {
			errCodes = append(errCodes, returnCodes...)
		}
	}

	if len(errCodes) == 0 {
		return nil
	}
	return errormgmt.CreateErrorResponseFromErrorCodesList(errCodes)
}

func (u UserBrImpl) validateUserNameNotTaken(dbctx context.Context, dto *userdtos.UserSaveDto) ([]string, error) {
	var errCodes []string
	if err := u.userRepository.FindByUsername(dbctx, dto.UserName, &models.User{}); err == nil {
		errCodes = append(errCodes, errormgmt.ErrCodeUsernameTaken)
	} else if !u.dbClient.IsNotFoundError(err) {
		return errCodes, err
	}
	return errCodes, nil
}

func NewUserBrImpl(dbClient dbaccess.DBClient, userRepository repositories.UserRepository) *UserBrImpl {
	return &UserBrImpl{dbClient: dbClient, userRepository: userRepository}
}
