package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/extensions/streamx/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
)

type UserBr interface {
	ValidateUserCreate(
		ctx context.Context,
		identity security.Identity,
		dto userdtos.UserSaveDto,
	) single.Single[[]apperrors.RuleError]

	ValidateUserUpdate(
		ctx context.Context,
		dto userdtos.UserSaveDto,
		existing models.User,
	) single.Single[[]apperrors.RuleError]
}

type UserBrImpl struct {
	dbHandler      database.DBHandler
	userRepository repositories.UserRepository
	errorService   errorservices.ErrorService
}

func (u UserBrImpl) ValidateUserCreate(
	ctx context.Context,
	identity security.Identity,
	dto userdtos.UserSaveDto,
) single.Single[[]apperrors.RuleError] {
	userNameNotTakenValidationSrc := u.validateUserNameNotTakenAsync(ctx, dto)
	userNotCreatedValidationSrc := validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByAuthIdAsync(ctx, identity.GetAuthId()),
		apperrors.ErrCodeUserAlreadyCreated,
	)
	ruleErrorsSrc := validationutils.ConcatSinglesOfRuleErrs(userNameNotTakenValidationSrc, userNotCreatedValidationSrc)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrorsSrc)
}

func (u UserBrImpl) ValidateUserUpdate(
	ctx context.Context,
	dto userdtos.UserSaveDto,
	existing models.User,
) single.Single[[]apperrors.RuleError] {
	ruleErrorsSrc := single.Just([]apperrors.RuleError{})
	if dto.UserName != existing.UserName {
		userNameNotTakenValidationSrc := u.validateUserNameNotTaken(ctx, dto)
		ruleErrorsSrc = validationutils.ConcatSinglesOfRuleErrs(ruleErrorsSrc, userNameNotTakenValidationSrc)
	}
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrorsSrc)
}

func (u UserBrImpl) validateUserNameNotTakenAsync(
	ctx context.Context,
	dto userdtos.UserSaveDto,
) single.Single[[]apperrors.RuleError] {
	return validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByUsernameAsync(ctx, dto.UserName),
		apperrors.ErrCodeUsernameTaken,
	)
}

func (u UserBrImpl) validateUserNameNotTaken(
	ctx context.Context,
	dto userdtos.UserSaveDto,
) single.Single[[]apperrors.RuleError] {
	return validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByUsername(ctx, dto.UserName),
		apperrors.ErrCodeUsernameTaken,
	)
}

func NewUserBrImpl(
	dbHandler database.DBHandler,
	userRepository repositories.UserRepository,
	errorMessageService errorservices.ErrorService,
) UserBr {
	return &UserBrImpl{dbHandler: dbHandler, userRepository: userRepository, errorService: errorMessageService}
}
