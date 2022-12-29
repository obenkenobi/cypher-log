package businessrules

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type NoteBr interface {
	ValidateNoteUpdate(userBo userbos.UserBo, keyDto keydtos.UserKeyDto, existing models.Note) error
	ValidateNoteRead(userBo userbos.UserBo, keyDto keydtos.UserKeyDto, existing models.Note) error
	ValidateNoteDelete(userBo userbos.UserBo, existing models.Note) error
	ValidateGetNotes(pageRequest pagination.PageRequest) error
}

type NoteBrImpl struct {
	errorService    sharedservices.ErrorService
	validSortFields map[string]any
}

func (n NoteBrImpl) ValidateGetNotes(pageRequest pagination.PageRequest) error {
	var ruleErrors []apperrors.RuleError
	if len(pageRequest.Sort) != 1 {
		ruleErrors = append(ruleErrors, n.errorService.RuleErrorFromCode(apperrors.ErrCodeMustSortByOneOption))
	} else if _, ok := n.validSortFields[pageRequest.Sort[0].Field]; !ok {
		ruleErrors = append(ruleErrors, n.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSortOptions))
	}
	return validationutils.MergeAppErrors(ruleErrors)
}

func (n NoteBrImpl) ValidateNoteRead(userBo userbos.UserBo, keyDto keydtos.UserKeyDto, existing models.Note) error {
	ruleErrs := append(n.validateKeyVersion(keyDto, existing), n.validateNoteOwnership(userBo, existing)...)
	return validationutils.MergeAppErrors(ruleErrs)
}

func (n NoteBrImpl) ValidateNoteUpdate(userBo userbos.UserBo, keyDto keydtos.UserKeyDto, existing models.Note) error {
	ruleErrs := append(n.validateKeyVersion(keyDto, existing), n.validateNoteOwnership(userBo, existing)...)
	return validationutils.MergeAppErrors(ruleErrs)
}

func (n NoteBrImpl) ValidateNoteDelete(userBo userbos.UserBo, existing models.Note) error {
	return validationutils.MergeAppErrors(n.validateNoteOwnership(userBo, existing))
}

func (n NoteBrImpl) validateKeyVersion(keyDto keydtos.UserKeyDto, existing models.Note) []apperrors.RuleError {
	var ruleErrs []apperrors.RuleError
	if keyDto.KeyVersion != existing.KeyVersion {
		ruleErrs = append(ruleErrs, n.errorService.RuleErrorFromCode(apperrors.ErrCodeDataRace))
	}
	return ruleErrs
}

func (n NoteBrImpl) validateNoteOwnership(userBo userbos.UserBo, existing models.Note) []apperrors.RuleError {
	var ruleErrs []apperrors.RuleError
	if userBo.Id != existing.UserId {
		ruleErrs = append(ruleErrs, n.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
	}
	return ruleErrs
}

func NewNoteBrImpl(errorService sharedservices.ErrorService) *NoteBrImpl {
	return &NoteBrImpl{
		errorService: errorService,
		validSortFields: map[string]any{
			pagination.SortFieldCreatedAt: any(true),
			pagination.SortFieldUpdatedAt: any(true),
		},
	}
}
