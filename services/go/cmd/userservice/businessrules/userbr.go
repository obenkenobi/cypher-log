package businessrules

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
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
	ruleErrorsX := stream.Just([]apperrors.RuleError{})

	ruleErrorsX = stream.FlatMap(
		ruleErrorsX,
		func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
			userFind := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
			return stream.Map(userFind, func(userMaybe option.Maybe[models.User]) []apperrors.RuleError {
				if userMaybe.IsPresent() {
					return append(ruleErrors, u.errorService.RuleErrorFromCode(apperrors.ErrCodeUserAlreadyCreated))
				}
				return ruleErrors
			})
		},
	)
	ruleErrorsX = stream.FlatMap(
		ruleErrorsX,
		func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
			userNameValidation := u.validateUserNameNotTaken(ctx, dto)
			return stream.Map(userNameValidation, func(userNameErrors []apperrors.RuleError) []apperrors.RuleError {
				return append(ruleErrors, userNameErrors...)
			})
		},
	)

	return stream.FlatMap(ruleErrorsX, func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
		if len(ruleErrors) == 0 {
			return stream.Just(ruleErrors)
		}
		return stream.Error[[]apperrors.RuleError](apperrors.NewBadReqErrorFromRuleErrors(ruleErrors...))
	})
}

func (u UserBrImpl) ValidateUserUpdate(
	ctx context.Context,
	dto userdtos.UserSaveDto,
	existing models.User,
) stream.Observable[[]apperrors.RuleError] {
	ruleErrorsX := stream.Just([]apperrors.RuleError{})
	if dto.UserName != existing.UserName {
		ruleErrorsX = stream.FlatMap(
			ruleErrorsX,
			func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
				userNameValidation := u.validateUserNameNotTaken(ctx, dto)
				return stream.Map(userNameValidation, func(userNameErrors []apperrors.RuleError) []apperrors.RuleError {
					return append(ruleErrors, userNameErrors...)
				})
			},
		)
	}
	return stream.FlatMap(ruleErrorsX, func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
		if len(ruleErrors) == 0 {
			return stream.Just(ruleErrors)
		}
		return stream.Error[[]apperrors.RuleError](apperrors.NewBadReqErrorFromRuleErrors(ruleErrors...))
	})
}

func (u UserBrImpl) validateUserNameNotTaken(
	ctx context.Context,
	dto userdtos.UserSaveDto,
) stream.Observable[[]apperrors.RuleError] {
	ruleErrorsX := stream.Just([]apperrors.RuleError{})
	ruleErrorsX = stream.FlatMap(
		ruleErrorsX,
		func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
			userFind := u.userRepository.FindByUsername(ctx, dto.UserName)
			return stream.Map(userFind, func(userMaybe option.Maybe[models.User]) []apperrors.RuleError {
				if userMaybe.IsPresent() {
					return append(ruleErrors, u.errorService.RuleErrorFromCode(apperrors.ErrCodeUsernameTaken))
				}
				return ruleErrors
			})
		},
	)
	return ruleErrorsX
}

func NewUserBrImpl(
	dbHandler database.DBHandler,
	userRepository repositories.UserRepository,
	errorMessageService apperrors.ErrorService,
) UserBr {
	return &UserBrImpl{dbHandler: dbHandler, userRepository: userRepository, errorService: errorMessageService}
}
