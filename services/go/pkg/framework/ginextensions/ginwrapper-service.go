package ginextensions

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GinWrapperService interface {
	readPathStr(c *gin.Context, name string) string
	HandleErrorResponse(c *gin.Context, err error)
	ProcessBindError(err error) errors.BadRequestError
	RespondJsonOk(c *gin.Context, model any, err error)
}

type GinWrapperServiceImpl struct {
	errorMessageService errors.ErrorService
}

func (h GinWrapperServiceImpl) readPathStr(c *gin.Context, key string) string {
	return c.Param(key)
}

func (h GinWrapperServiceImpl) ProcessBindError(err error) errors.BadRequestError {
	if verrors, ok := err.(validator.ValidationErrors); ok {
		var validationErrors []errors.ValidationError
		for _, verr := range verrors {
			validationError := errors.ValidationError{Field: verr.Field(), Message: verr.ActualTag()}
			validationErrors = append(validationErrors, validationError)
		}
		return errors.NewBadReqErrorFromValidationErrors(validationErrors)
	}
	log.WithError(err).Info("Unable to bind json")
	return errors.NewBadReqErrorFromRuleError(h.errorMessageService.RuleErrorFromCode(errors.ErrCodeCannotBindJson))
}

func (h GinWrapperServiceImpl) HandleErrorResponse(c *gin.Context, err error) {
	if badReqErr, ok := err.(errors.BadRequestError); ok {
		c.JSON(http.StatusBadRequest, badReqErr)
	} else {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
	}
}

func (h GinWrapperServiceImpl) RespondJsonOk(c *gin.Context, model any, err error) {
	if err != nil {
		h.HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, model)
}

func NewGinWrapperService(errorService errors.ErrorService) GinWrapperService {
	return &GinWrapperServiceImpl{errorMessageService: errorService}
}

func BindBody[V any](ginWrapperService GinWrapperService, c *gin.Context, obj V) stream.Observable[V] {
	if err := c.ShouldBind(obj); err != nil {
		return stream.Error[V](ginWrapperService.ProcessBindError(err))
	}
	return stream.Just(obj)
}
