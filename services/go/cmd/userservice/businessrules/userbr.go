package businessrules

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
)

type UserBr interface {
	ValidateUserCreate(
		ctx context.Context,
		identity security.Identity,
		dto *userdtos.UserSaveDto,
	) stream.Observable[[]errors.RuleError]

	ValidateUserUpdate(
		ctx context.Context,
		dto *userdtos.UserSaveDto,
		existing *models.User,
	) stream.Observable[[]errors.RuleError]
}

type UserBrImpl struct {
	dbHandler      database.DBHandler
	userRepository repositories.UserRepository
	errorService   errors.ErrorService
}

func (u UserBrImpl) ValidateUserCreate(
	ctx context.Context,
	identity security.Identity,
	dto *userdtos.UserSaveDto,
) stream.Observable[[]errors.RuleError] {
	ruleErrorsX := stream.Just([]errors.RuleError{})

	ruleErrorsX = stream.FlatMap(
		ruleErrorsX,
		func(ruleErrors []errors.RuleError) stream.Observable[[]errors.RuleError] {
			userFind := u.userRepository.FindByAuthId(ctx, identity.GetAuthId())
			return stream.Map(userFind, func(userMaybe option.Maybe[*models.User]) []errors.RuleError {
				if userMaybe.IsPresent() {
					return append(ruleErrors, u.errorService.RuleErrorFromCode(errors.ErrCodeUserAlreadyCreated))
				}
				return ruleErrors
			})
		},
	)
	ruleErrorsX = stream.FlatMap(
		ruleErrorsX,
		func(ruleErrors []errors.RuleError) stream.Observable[[]errors.RuleError] {
			userNameValidation := u.validateUserNameNotTaken(ctx, dto)
			return stream.Map(userNameValidation, func(unameruleErrors []errors.RuleError) []errors.RuleError {
				return append(ruleErrors, unameruleErrors...)
			})
		},
	)

	return stream.FlatMap(ruleErrorsX, func(ruleErrors []errors.RuleError) stream.Observable[[]errors.RuleError] {
		if len(ruleErrors) == 0 {
			return stream.Just(ruleErrors)
		}
		return stream.Error[[]errors.RuleError](errors.NewBadReqErrorFromRuleErrors(ruleErrors...))
	})
}

func (u UserBrImpl) ValidateUserUpdate(
	ctx context.Context,
	dto *userdtos.UserSaveDto,
	existing *models.User,
) stream.Observable[[]errors.RuleError] {
	ruleErrorsX := stream.Just([]errors.RuleError{})
	if dto.UserName != existing.UserName {
		ruleErrorsX = stream.FlatMap(
			ruleErrorsX,
			func(ruleErrors []errors.RuleError) stream.Observable[[]errors.RuleError] {
				userNameValidation := u.validateUserNameNotTaken(ctx, dto)
				return stream.Map(userNameValidation, func(unameruleErrors []errors.RuleError) []errors.RuleError {
					return append(ruleErrors, unameruleErrors...)
				})
			},
		)
	}
	return stream.FlatMap(ruleErrorsX, func(ruleErrors []errors.RuleError) stream.Observable[[]errors.RuleError] {
		if len(ruleErrors) == 0 {
			return stream.Just(ruleErrors)
		}
		return stream.Error[[]errors.RuleError](errors.NewBadReqErrorFromRuleErrors(ruleErrors...))
	})
}

func (u UserBrImpl) validateUserNameNotTaken(
	ctx context.Context,
	dto *userdtos.UserSaveDto,
) stream.Observable[[]errors.RuleError] {
	ruleErrorsX := stream.Just([]errors.RuleError{})
	ruleErrorsX = stream.FlatMap(
		ruleErrorsX,
		func(ruleErrors []errors.RuleError) stream.Observable[[]errors.RuleError] {
			userFind := u.userRepository.FindByUsername(ctx, dto.UserName)
			return stream.Map(userFind, func(userMaybe option.Maybe[*models.User]) []errors.RuleError {
				if userMaybe.IsPresent() {
					return append(ruleErrors, u.errorService.RuleErrorFromCode(errors.ErrCodeUsernameTaken))
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
	errorMessageService errors.ErrorService,
) UserBr {
	return &UserBrImpl{dbHandler: dbHandler, userRepository: userRepository, errorService: errorMessageService}
}
