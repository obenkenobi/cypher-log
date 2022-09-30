package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
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
	crudDBHandler  dshandlers.CrudDSHandler
	userRepository repositories.UserRepository
	errorService   sharedservices.ErrorService
}

func (u UserBrImpl) ValidateUserCreate(
	ctx context.Context,
	identity security.Identity,
	dto userdtos.UserSaveDto,
) single.Single[[]apperrors.RuleError] {
	userNameNotTakenValidationSrc := u.validateUserNameNotTaken(ctx, dto).ScheduleEagerAsync(ctx)
	userNotCreatedValidationSrc := validationutils.ValidateValueIsNotPresent(
		u.errorService,
		u.userRepository.FindByAuthId(ctx, identity.GetAuthId()).ScheduleEagerAsync(ctx),
		apperrors.ErrCodeResourceAlreadyCreated,
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
	crudDBHandler dshandlers.CrudDSHandler,
	userRepository repositories.UserRepository,
	errorMessageService sharedservices.ErrorService,
) UserBr {
	return &UserBrImpl{crudDBHandler: crudDBHandler, userRepository: userRepository, errorService: errorMessageService}
}
