package validationutils

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

func ValidateValueIsNotPresent[V any](
	errorService sharedservices.ErrorService,
	maybe option.Maybe[V],
	notPresentErrorCode string,
) []apperrors.RuleError {
	if maybe.IsPresent() {
		return []apperrors.RuleError{errorService.RuleErrorFromCode(notPresentErrorCode)}
	}
	return []apperrors.RuleError{}
}

func ValidateValueIsPresent[V any](
	errorService sharedservices.ErrorService,
	maybe option.Maybe[V],
	notPresentErrorCode string,
) []apperrors.RuleError {
	if maybe.IsEmpty() {
		return []apperrors.RuleError{errorService.RuleErrorFromCode(notPresentErrorCode)}
	}
	return []apperrors.RuleError{}
}

func MergeRuleErrors(ruleErrs []apperrors.RuleError) error {
	if len(ruleErrs) == 0 {
		return nil
	}
	return apperrors.NewBadReqErrorFromRuleErrors(ruleErrs...)
}
