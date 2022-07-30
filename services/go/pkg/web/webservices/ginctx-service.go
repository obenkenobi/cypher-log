package webservices

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors/errorservices"
	"github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GinCtxService interface {
	readPathStr(c *gin.Context, name string) string
	HandleErrorResponse(c *gin.Context, err error)
	ProcessBindError(err error) apperrors.BadRequestError
	RespondJsonOk(c *gin.Context, model any, err error)
}

type GinCtxServiceImpl struct {
	errorMessageService errorservices.ErrorService
}

func (h GinCtxServiceImpl) readPathStr(c *gin.Context, key string) string {
	return c.Param(key)
}

func (h GinCtxServiceImpl) ProcessBindError(err error) apperrors.BadRequestError {
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

func (h GinCtxServiceImpl) HandleErrorResponse(c *gin.Context, err error) {
	if badReqErr, ok := err.(apperrors.BadRequestError); ok {
		c.JSON(http.StatusBadRequest, badReqErr)
	} else {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
	}
}

func (h GinCtxServiceImpl) RespondJsonOk(c *gin.Context, model any, err error) {
	if err != nil {
		h.HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, model)
}

func NewGinWrapperService(errorService errorservices.ErrorService) GinCtxService {
	return &GinCtxServiceImpl{errorMessageService: errorService}
}

func BindValueToBody[V any](ginCtxService GinCtxService, c *gin.Context, value V) single.Single[V] {
	if err := c.ShouldBind(&value); err != nil {
		return single.Error[V](ginCtxService.ProcessBindError(err))
	}
	return single.Just(value)
}

func BindPointerToBody[V any](ginCtxService GinCtxService, c *gin.Context, value *V) single.Single[*V] {
	if err := c.ShouldBind(value); err != nil {
		return single.Error[*V](ginCtxService.ProcessBindError(err))
	}
	return single.Just(value)
}
