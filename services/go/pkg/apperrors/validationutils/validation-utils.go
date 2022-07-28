package validationutils

import (
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
)

func ValidateValueIsNotPresent[V any](
	errorService apperrors.ErrorService,
	observableX stream.Observable[option.Maybe[V]],
	notPresentErrorCode string,
) stream.Observable[[]apperrors.RuleError] {
	return stream.Map(observableX, func(maybe option.Maybe[V]) []apperrors.RuleError {
		if maybe.IsPresent() {
			return []apperrors.RuleError{errorService.RuleErrorFromCode(notPresentErrorCode)}
		}
		return []apperrors.RuleError{}
	})
}

func ValidateValueIsPresent[V any](
	errorService apperrors.ErrorService,
	observableX stream.Observable[option.Maybe[V]],
	notPresentErrorCode string,
) stream.Observable[[]apperrors.RuleError] {
	return stream.Map(observableX, func(maybe option.Maybe[V]) []apperrors.RuleError {
		if maybe.IsEmpty() {
			return []apperrors.RuleError{errorService.RuleErrorFromCode(notPresentErrorCode)}
		}
		return []apperrors.RuleError{}
	})
}

func ConcatRuleErrorObservables(
	src1 stream.Observable[[]apperrors.RuleError],
	src2 stream.Observable[[]apperrors.RuleError],
) stream.Observable[[]apperrors.RuleError] {
	return stream.FlatMap(src1, func(srcErrs []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
		return stream.Map(src2, func(ruleErrs []apperrors.RuleError) []apperrors.RuleError {
			return append(ruleErrs, srcErrs...)
		})
	})
	//errorsX := stream.Reduce(
	//	stream.Concat(ruleErrorObservables...),
	//	[]apperrors.RuleError{},
	//	func(ruleErrors1, ruleErrors2 []apperrors.RuleError) []apperrors.RuleError {
	//		return append(ruleErrors1, ruleErrors2...)
	//	},
	//)
	//return stream.Take(1, errorsX)
}

func PassRuleErrorsIfEmptyElsePassBadReqError(
	ruleErrorsX stream.Observable[[]apperrors.RuleError],
) stream.Observable[[]apperrors.RuleError] {
	return stream.FlatMap(ruleErrorsX, func(ruleErrors []apperrors.RuleError) stream.Observable[[]apperrors.RuleError] {
		if len(ruleErrors) == 0 {
			return stream.Just(ruleErrors)
		}
		return stream.Error[[]apperrors.RuleError](apperrors.NewBadReqErrorFromRuleErrors(ruleErrors...))
	})
}
