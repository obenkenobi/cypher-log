package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type UserBr interface {
	ValidateUserCreate(ctx context.Context, identity security.Identity, dto userdtos.UserSaveDto) error
	ValidateUserUpdate(ctx context.Context, dto userdtos.UserSaveDto, existing models.User) error
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
) error {
	userNameNotTakenValidationErrs, err := u.validateUserNameNotTaken(ctx, dto)
	if err != nil {
		return err
	}
	userFindByAuthIdMaybe, err := u.userRepository.FindByAuthIdAndNotToBeDeleted(ctx, identity.GetAuthId())
	if err != nil {
		return err
	}
	userNotCreatedValidationErrs := validationutils.ValidateValueIsNotPresent(
		u.errorService,
		userFindByAuthIdMaybe,
		apperrors.ErrCodeResourceAlreadyCreated,
	)
	ruleErrorsSrc := append(userNameNotTakenValidationErrs, userNotCreatedValidationErrs...)
	return validationutils.MergeAppErrors(ruleErrorsSrc)
}

func (u UserBrImpl) ValidateUserUpdate(ctx context.Context, dto userdtos.UserSaveDto, existing models.User) error {
	var ruleErrs []apperrors.RuleError
	if dto.UserName != existing.UserName {
		valErrs, err := u.validateUserNameNotTaken(ctx, dto)
		if err != nil {
			return err
		}
		ruleErrs = append(ruleErrs, valErrs...)
	}
	return validationutils.MergeAppErrors(ruleErrs)
}

func (u UserBrImpl) validateUserNameNotTaken(
	ctx context.Context,
	dto userdtos.UserSaveDto,
) ([]apperrors.RuleError, error) {
	maybe, err := u.userRepository.FindByUsernameAndNotToBeDeleted(ctx, dto.UserName)
	if err != nil {
		return nil, err
	}
	return validationutils.ValidateValueIsNotPresent(u.errorService, maybe, apperrors.ErrCodeUsernameTaken), nil
}

func NewUserBrImpl(
	crudDBHandler dshandlers.CrudDSHandler,
	userRepository repositories.UserRepository,
	errorMessageService sharedservices.ErrorService,
) *UserBrImpl {
	return &UserBrImpl{crudDBHandler: crudDBHandler, userRepository: userRepository, errorService: errorMessageService}
}
