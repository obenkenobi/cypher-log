package errors

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
)

func parseValidationError(verrors validator.ValidationErrors) *errordtos.ErrorResponseDto {
	var validationErrorDtos []errordtos.ValidationErrDto
	for _, validationErr := range verrors {
		validationErrorDtos = append(validationErrorDtos, errordtos.ValidationErrDto{
			Field:   validationErr.Field(),
			Message: validationErr.ActualTag(),
		})
	}
	return errordtos.NewValidationErrorResponse(validationErrorDtos)
}

func HandleBindError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		HandleErrorResponse(c, *parseValidationError(validationErrors))
		return
	}
	// We now know that this error is not a validation error
	// probably a malformed JSON
	log.WithError(err).Info("Unable to bind")
	HandleErrorResponse(c, *errordtos.NewErrorResponse("Invalid input", false))
}

func HandleErrorResponse(c *gin.Context, errorResponse errordtos.ErrorResponseDto) {
	var httpStatus int
	if errorResponse.IsInternalError {
		httpStatus = http.StatusInternalServerError
	} else {
		httpStatus = http.StatusBadRequest
	}
	c.JSON(httpStatus, errorResponse)
}
