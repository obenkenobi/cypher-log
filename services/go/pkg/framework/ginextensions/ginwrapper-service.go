package ginextensions

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GinWrapperService interface {
	readPathStr(c *gin.Context, name string) string
	HandleErrorResponse(c *gin.Context, err error)
	ProcessBindError(err error) apperrors.BadRequestError
	RespondJsonOk(c *gin.Context, model any, err error)
}

type GinWrapperServiceImpl struct {
	errorMessageService apperrors.ErrorService
}

func (h GinWrapperServiceImpl) readPathStr(c *gin.Context, key string) string {
	return c.Param(key)
}

func (h GinWrapperServiceImpl) ProcessBindError(err error) apperrors.BadRequestError {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		var appValErrors []apperrors.ValidationError
		for _, fieldError := range fieldErrors {
			appValError := apperrors.ValidationError{Field: fieldError.Field(), Message: fieldError.ActualTag()}
			appValErrors = append(appValErrors, appValError)
		}
		return apperrors.NewBadReqErrorFromValidationErrors(appValErrors)
	}
	log.WithError(err).Info("Unable to bind json")
	return apperrors.NewBadReqErrorFromRuleError(h.errorMessageService.RuleErrorFromCode(apperrors.ErrCodeCannotBindJson))
}

func (h GinWrapperServiceImpl) HandleErrorResponse(c *gin.Context, err error) {
	if badReqErr, ok := err.(apperrors.BadRequestError); ok {
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

func NewGinWrapperService(errorService apperrors.ErrorService) GinWrapperService {
	return &GinWrapperServiceImpl{errorMessageService: errorService}
}

func BindValueToBody[V any](ginWrapperService GinWrapperService, c *gin.Context, value V) stream.Observable[V] {
	if err := c.ShouldBind(&value); err != nil {
		return stream.Error[V](ginWrapperService.ProcessBindError(err))
	}
	return stream.Just(value)
}

func BindPointerToBody[V any](ginWrapperService GinWrapperService, c *gin.Context, value *V) stream.Observable[*V] {
	if err := c.ShouldBind(value); err != nil {
		return stream.Error[*V](ginWrapperService.ProcessBindError(err))
	}
	return stream.Just(value)
}
