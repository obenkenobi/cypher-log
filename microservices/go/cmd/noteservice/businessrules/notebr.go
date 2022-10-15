package businessrules

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
)

type NoteBr interface {
	ValidateNoteUpdate(
		userBo userbos.UserBo,
		session commondtos.UKeySessionDto,
		existing models.Note,
	) single.Single[[]apperrors.RuleError]
}

type NoteBrImpl struct {
	errorService sharedservices.ErrorService
}

func (n NoteBrImpl) ValidateNoteUpdate(userBo userbos.UserBo, session commondtos.UKeySessionDto, existing models.Note) single.Single[[]apperrors.RuleError] {
	validateKeyVersionSrc := n.validateKeyVersion(session, existing)
	validateNoteOwnedByUser := single.FromSupplierCached(func() ([]apperrors.RuleError, error) {
		var ruleErrs []apperrors.RuleError
		if userBo.Id != existing.UserId {
			ruleErrs = append(ruleErrs, n.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound))
		}
		return ruleErrs, nil
	})
	ruleErrs := validationutils.ConcatSinglesOfRuleErrs(validateKeyVersionSrc, validateNoteOwnedByUser)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrs)
}

func (n NoteBrImpl) validateKeyVersion(session commondtos.UKeySessionDto, existing models.Note) single.Single[[]apperrors.RuleError] {
	return single.FromSupplierCached(func() ([]apperrors.RuleError, error) {
		var ruleErrs []apperrors.RuleError
		if session.KeyVersion != existing.KeyVersion {
			ruleErrs = append(ruleErrs, n.errorService.RuleErrorFromCode(apperrors.ErrCodeDataRace))
		}
		return ruleErrs, nil
	})
}

func NewNoteBrImpl(errorService sharedservices.ErrorService) *NoteBrImpl {
	return &NoteBrImpl{errorService: errorService}
}
