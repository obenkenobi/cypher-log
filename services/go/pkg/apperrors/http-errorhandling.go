package apperrors

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
)

func HandleBindError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		HandleErrorResponse(c, *CreateErrorResponseFromValidationErrors(validationErrors))
		return
	}
	// We now know that this error is not a validation error
	// probably a malformed JSON
	log.WithError(err).Info("Unable to bind")

	HandleErrorResponse(c, *CreateErrorResponseFromErrorCodes(ErrCodeCannotBindJson))
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
