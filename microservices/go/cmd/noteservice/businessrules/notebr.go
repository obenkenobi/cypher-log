package businessrules

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
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
	ValidateGetNotes(pageRequest pagination.PageRequest) single.Single[any]
}

type NoteBrImpl struct {
	errorService    sharedservices.ErrorService
	validSortFields map[string]any
}

func (n NoteBrImpl) ValidateGetNotes(pageRequest pagination.PageRequest) single.Single[any] {
	validateSortingSrc := single.FromSupplierCached(func() ([]apperrors.RuleError, error) {
		var ruleErrors []apperrors.RuleError
		if len(pageRequest.Sort) != 1 {
			ruleErrors = append(ruleErrors, n.errorService.RuleErrorFromCode(apperrors.ErrCodeMustSortByOneOption))
		} else if _, ok := n.validSortFields[pageRequest.Sort[0].Field]; !ok {
			ruleErrors = append(ruleErrors, n.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSortOptions))
		}
		return ruleErrors, nil
	})
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(validateSortingSrc)
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
	return &NoteBrImpl{
		errorService: errorService,
		validSortFields: map[string]any{
			"created_at": any(true),
			"updated_at": any(true),
		},
	}
}
