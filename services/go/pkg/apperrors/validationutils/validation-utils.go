package validationutils

import (
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/services/go/pkg/extensions/streamx/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
)

func ValidateValueIsNotPresent[V any](
	errorService errorservices.ErrorService,
	valSrc single.Single[option.Maybe[V]],
	notPresentErrorCode string,
) single.Single[[]apperrors.RuleError] {
	return single.Map(valSrc, func(maybe option.Maybe[V]) []apperrors.RuleError {
		if maybe.IsPresent() {
			return []apperrors.RuleError{errorService.RuleErrorFromCode(notPresentErrorCode)}
		}
		return []apperrors.RuleError{}
	})
}

func ValidateValueIsPresent[V any](
	errorService errorservices.ErrorService,
	valSrc single.Single[option.Maybe[V]],
	notPresentErrorCode string,
) single.Single[[]apperrors.RuleError] {
	return single.Map(valSrc, func(maybe option.Maybe[V]) []apperrors.RuleError {
		if maybe.IsEmpty() {
			return []apperrors.RuleError{errorService.RuleErrorFromCode(notPresentErrorCode)}
		}
		return []apperrors.RuleError{}
	})
}

func ConcatSinglesOfRuleErrs(
	src1 single.Single[[]apperrors.RuleError],
	src2 single.Single[[]apperrors.RuleError],
) single.Single[[]apperrors.RuleError] {
	return single.Map(
		single.Zip2(src1, src2),
		func(rulErrsTuple tuple.T2[[]apperrors.RuleError, []apperrors.RuleError]) []apperrors.RuleError {
			return append(rulErrsTuple.V1, rulErrsTuple.V2...)
		},
	)
}

func PassRuleErrorsIfEmptyElsePassBadReqError(
	ruleErrsSrc single.Single[[]apperrors.RuleError],
) single.Single[[]apperrors.RuleError] {
	return single.FlatMap(ruleErrsSrc, func(ruleErrors []apperrors.RuleError) single.Single[[]apperrors.RuleError] {
		if len(ruleErrors) == 0 {
			return single.Just(ruleErrors)
		}
		return single.Error[[]apperrors.RuleError](apperrors.NewBadReqErrorFromRuleErrors(ruleErrors...))
	})
}