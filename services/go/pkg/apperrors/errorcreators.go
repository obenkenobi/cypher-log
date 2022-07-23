package apperrors

import (
	"github.com/go-playground/validator/v10"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	log "github.com/sirupsen/logrus"
)

func getMsgFromCode(code string) string {
	if msg, ok := GetErrorCodeToMsgMap()[code]; ok {
		return msg
	}
	return code
}

func GetAppErrorDtosFromCodes(codes []string) []errordtos.AppErrorDto {
	codeToMsgMap := map[string]string{}
	for _, code := range codes {
		codeToMsgMap[code] = getMsgFromCode(code)
	}
	var appErrors []errordtos.AppErrorDto
	for code, msg := range codeToMsgMap {
		appErrors = append(appErrors, errordtos.AppErrorDto{Code: code, Message: msg})
	}
	return appErrors
}

func CreateErrorResponseFromErrorCodes(codes ...string) *errordtos.ErrorResponseDto {
	return CreateErrorResponseFromAppErrorDtos(GetAppErrorDtosFromCodes(codes)...)
}

func CreateErrorResponseFromAppErrorDtos(appErrorDtos ...errordtos.AppErrorDto) *errordtos.ErrorResponseDto {
	return errordtos.NewAppErrorsResponse(appErrorDtos)
}

func CreateErrorResponseFromValidationErrors(verrors validator.ValidationErrors) *errordtos.ErrorResponseDto {
	var validationErrorDtos []errordtos.ValidationErrDto
	for _, validationErr := range verrors {
		validationErrorDtos = append(validationErrorDtos, errordtos.ValidationErrDto{
			Field:   validationErr.Field(),
			Message: validationErr.ActualTag(),
		})
	}
	return errordtos.NewValidationErrorResponse(validationErrorDtos)
}

func CreateInternalErrResponse(err error) *errordtos.ErrorResponseDto {
	log.WithError(err).Error()
	return errordtos.NewInternalErrResponse()
}
