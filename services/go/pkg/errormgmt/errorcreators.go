package errormgmt

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/sirupsen/logrus"
)

func CreateErrorResponseFromErrorCodes(codes ...string) *errordtos.ErrorResponseDto {
	return CreateErrorResponseFromErrorCodesList(codes)
}

func CreateErrorResponseFromErrorCodesList(codes []string) *errordtos.ErrorResponseDto {
	return errordtos.NewAppErrorsResponse(GetAppErrorDtosFromCodes(codes))
}

func CreateInternalErrResponseWithErrLog(err error) *errordtos.ErrorResponseDto {
	logrus.WithError(err).Error()
	return errordtos.NewInternalErrResponse()
}
