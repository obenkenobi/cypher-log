package grpcerrorhandling

import (
	"encoding/json"
	"github.com/akrennmair/slice"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ProcessError(err error) error {
	if err == nil {
		return err
	}
	if badReqErr, ok := err.(apperrors.BadRequestError); ok {
		var badReqJson string
		if jsonBytes, err := json.Marshal(badReqErr); err != nil {
			badReqJson = "{}"
		} else {
			badReqJson = string(jsonBytes)
		}
		if len(badReqErr.RuleErrors) > 1 {
			return status.Error(codes.FailedPrecondition, badReqJson)
		} else if len(badReqErr.RuleErrors) == 1 {
			type violations map[string]apperrors.RuleError
			brViolations := slice.Reduce(badReqErr.RuleErrors,
				func(m violations, ruleErr apperrors.RuleError) violations { m[ruleErr.Code] = ruleErr; return m })
			if _, ok := brViolations[apperrors.ErrCodeReqResourcesNotFound]; ok {
				return status.Error(codes.NotFound, badReqJson)
			} else if _, ok := brViolations[apperrors.ErrCodeResourceAlreadyCreated]; ok {
				return status.Error(codes.AlreadyExists, badReqJson)
			} else {
				return status.Error(codes.FailedPrecondition, badReqJson)
			}
		} else if len(badReqErr.ValidationErrors) > 0 {
			return status.Error(codes.InvalidArgument, badReqJson)
		} else {
			return status.Error(codes.Unknown, badReqJson)
		}
	}
	log.Error(err)
	return status.Error(codes.Internal, "internal error")
}
