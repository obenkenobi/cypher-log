package businessrules

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type NoteBr interface {
	ValidateNoteUpdate(
		userBo userbos.UserBo,
		keyDto keydtos.UserKeyDto,
		existing models.Note,
	) single.Single[any]
	ValidateNoteRead(userBo userbos.UserBo, keyDto keydtos.UserKeyDto, existing models.Note) single.Single[any]
	ValidateNoteDelete(userBo userbos.UserBo, existing models.Note) single.Single[any]
}

type NoteBrImpl struct {
	errorService sharedservices.ErrorService
}

func (n NoteBrImpl) ValidateNoteRead(
	userBo userbos.UserBo,
	keyDto keydtos.UserKeyDto,
	existing models.Note,
) single.Single[any] {
	validateKeyVersionSrc := n.validateKeyVersion(keyDto, existing)
	validateNoteOwnedByUser := n.validateNoteOwnership(userBo, existing)
	ruleErrs := validationutils.ConcatSinglesOfRuleErrs(validateKeyVersionSrc, validateNoteOwnedByUser)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrs)
}

func (n NoteBrImpl) ValidateNoteUpdate(
	userBo userbos.UserBo,
	keyDto keydtos.UserKeyDto,
	existing models.Note,
) single.Single[any] {
	validateKeyVersionSrc := n.validateKeyVersion(keyDto, existing)
	validateNoteOwnedByUser := n.validateNoteOwnership(userBo, existing)
	ruleErrs := validationutils.ConcatSinglesOfRuleErrs(validateKeyVersionSrc, validateNoteOwnedByUser)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrs)
}

func (n NoteBrImpl) ValidateNoteDelete(userBo userbos.UserBo, existing models.Note) single.Single[any] {
	validateNoteOwnedByUser := n.validateNoteOwnership(userBo, existing)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(validateNoteOwnedByUser)
}

func (n NoteBrImpl) validateKeyVersion(
	keyDto keydtos.UserKeyDto,
	existing models.Note,
) single.Single[[]apperrors.RuleError] {
	return single.FromSupplierCached(func() ([]apperrors.RuleError, error) {
		var ruleErrs []apperrors.RuleError
		if keyDto.KeyVersion != existing.KeyVersion {
			ruleErrs = append(ruleErrs, n.errorService.RuleErrorFromCode(apperrors.ErrCodeDataRace))
		}
		return ruleErrs, nil
	})
}

func (n NoteBrImpl) validateNoteOwnership(
	userBo userbos.UserBo,
	existing models.Note,
) single.Single[[]apperrors.RuleError] {
	return single.FromSupplierCached(func() ([]apperrors.RuleError, error) {
		var ruleErrs []apperrors.RuleError
		if userBo.Id != existing.UserId {
			ruleErrs = append(ruleErrs, n.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
		}
		return ruleErrs, nil
	})
}

func NewNoteBrImpl(errorService sharedservices.ErrorService) *NoteBrImpl {
	return &NoteBrImpl{errorService: errorService}
}
