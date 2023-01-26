package gtools

import (
	"context"
	"encoding/json"
	"github.com/akrennmair/slice"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type violationsMap map[string]apperrors.RuleError

func (v violationsMap) containsCode(code string) bool {
	_, ok := v[code]
	return ok
}

func badReqErrToJsonString(badReqErr apperrors.BadRequestError) string {
	if jsonBytes, err := json.Marshal(badReqErr); err == nil {
		return string(jsonBytes)
	}
	return "{}"
}

func processBadReqErr(procedureType GrpcAction, badReqErr apperrors.BadRequestError) error {
	var badReqJson string
	if bytes, err := json.Marshal(&badReqErr); err != nil {
		badReqJson = "{}"
	} else {
		badReqJson = string(bytes)
	}
	brViolations := slice.ReduceWithInitialValue(badReqErr.RuleErrors, make(violationsMap),
		func(m violationsMap, ruleErr apperrors.RuleError) violationsMap { m[ruleErr.Code] = ruleErr; return m })
	if _, ok := brViolations[apperrors.ErrCodeReqResourcesNotFound]; procedureType == ReadAction && ok {
		return status.Error(codes.NotFound, badReqJson)
	} else if _, ok := brViolations[apperrors.ErrCodeResourceAlreadyCreated]; procedureType == CreateAction && ok {
		return status.Error(codes.AlreadyExists, badReqJson)
	} else if len(badReqErr.ValidationErrors) > 0 {
		return status.Error(codes.InvalidArgument, badReqJson)
	} else if len(badReqErr.RuleErrors) > 0 {
		return status.Error(codes.FailedPrecondition, badReqJson)
	} else {
		return status.Error(codes.Unknown, badReqJson)
	}
}

// ProcessErrorToGrpcStatusError takes an error and procedure type and transforms the
// error into a format used by GRPC methods.
func ProcessErrorToGrpcStatusError(ctx context.Context, procedureType GrpcAction, err error) error {
	if err == nil {
		return err
	}
	if badReqErr, ok := err.(apperrors.BadRequestError); ok {
		return processBadReqErr(procedureType, badReqErr)
	}
	logger.Log.WithContext(ctx).Error(err)
	return status.Error(codes.Internal, "internal error")
}

type ErrorResponseHandler struct {
	err         error
	status      *status.Status
	badReqError *apperrors.BadRequestError
}

func (h ErrorResponseHandler) IfBadRequestError(handleFunc func(apperrors.BadRequestError)) ErrorResponseHandler {
	if h.badReqError != nil {
		handleFunc(*h.badReqError)
	}
	return h
}

func (h ErrorResponseHandler) IfOtherError(handleFunc func(error, *status.Status)) ErrorResponseHandler {
	if h.badReqError == nil {
		st := h.status
		if st == nil || st.Err() == nil {
			st = status.New(codes.Unknown, h.err.Error())
		}
		handleFunc(h.err, st)
	}
	return h
}

func (h ErrorResponseHandler) GetProcessedError() error {
	if h.badReqError != nil {
		return *h.badReqError
	}
	return h.err
}

func NewErrorResponseHandler(err error) ErrorResponseHandler {
	handler := ErrorResponseHandler{err: nil}
	if err == nil {
		return handler
	}
	handler.err = err
	st := status.Convert(err)
	handler.status = st
	var badReqErr apperrors.BadRequestError
	if err := json.Unmarshal([]byte(st.Message()), &badReqErr); err != nil {
		return handler
	}
	handler.badReqError = &badReqErr
	return handler
}
