package businessrules

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
)

type UserBr interface {
	ValidateUserCreate(
		ctx context.Context,
		identity security.Identity,
		dto userdtos.UserSaveDto,
	) stream.Observable[[]apperrors.RuleError]

	ValidateUserUpdate(
		ctx context.Context,
		dto userdtos.UserSaveDto,
		existing models.User,
	) stream.Observable[[]apperrors.RuleError]
}

type UserBrImpl struct {
	dbHandler      database.DBHandler
	userRepository repositories.UserRepository
	errorService   apperrors.ErrorService
}

func (u UserBrImpl) ValidateUserCreate(
	ctx context.Context,
	identity security.Identity,
	dto userdtos.UserSaveDto,
) stream.Observable[[]apperrors.RuleError] {
	userNameNotTakenValidationX := u.validateUserNameNotTakenAsync(ctx, dto)
	userNotCreatedValidationX := validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByAuthIdAsync(ctx, identity.GetAuthId()),
		apperrors.ErrCodeUserAlreadyCreated,
	)
	ruleErrorsX := validationutils.ConcatRuleErrorObservables(userNameNotTakenValidationX, userNotCreatedValidationX)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrorsX)
}

func (u UserBrImpl) ValidateUserUpdate(
	ctx context.Context,
	dto userdtos.UserSaveDto,
	existing models.User,
) stream.Observable[[]apperrors.RuleError] {
	ruleErrorsX := stream.Just([]apperrors.RuleError{})
	if dto.UserName != existing.UserName {
		serNameNotTakenValidationX := u.validateUserNameNotTaken(ctx, dto)
		ruleErrorsX = validationutils.ConcatRuleErrorObservables(ruleErrorsX, serNameNotTakenValidationX)
	}
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrorsX)
}

func (u UserBrImpl) validateUserNameNotTakenAsync(
	ctx context.Context,
	dto userdtos.UserSaveDto,
) stream.Observable[[]apperrors.RuleError] {
	return validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByUsernameAsync(ctx, dto.UserName),
		apperrors.ErrCodeUsernameTaken,
	)
}

func (u UserBrImpl) validateUserNameNotTaken(
	ctx context.Context,
	dto userdtos.UserSaveDto,
) stream.Observable[[]apperrors.RuleError] {
	return validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByUsername(ctx, dto.UserName),
		apperrors.ErrCodeUsernameTaken,
	)
}

func NewUserBrImpl(
	dbHandler database.DBHandler,
	userRepository repositories.UserRepository,
	errorMessageService apperrors.ErrorService,
) UserBr {
	return &UserBrImpl{dbHandler: dbHandler, userRepository: userRepository, errorService: errorMessageService}
}
