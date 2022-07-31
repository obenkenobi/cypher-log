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

// GinCtxService is an injectable service that provides convenience methods
// relating to http requests and responses using the Gin framework. As the name
// suggest, it is designed to handle a Gin's Context struct.
type GinCtxService interface {
	// ParamStr returns the value of a URL param as a string
	ParamStr(c *gin.Context, name string) string

	// HandleErrorResponse takes an error and parses it to set the appropriate http
	// response. Certain errors relating to user input will trigger a 4XX status
	// code. Otherwise, a 5XX code will be thrown indicating the error means
	// something went wrong with the server.
	HandleErrorResponse(c *gin.Context, err error)

	// processBindError takes an error from binding a value from a request body processes it into a BadRequestError.
	processBindError(err error) apperrors.BadRequestError

	// RespondJsonOk responds with json value with a 200 status code.
	RespondJsonOk(c *gin.Context, model any, err error)
}

type GinCtxServiceImpl struct {
	errorMessageService errorservices.ErrorService
}

func (h GinCtxServiceImpl) ParamStr(c *gin.Context, key string) string {
	return c.Param(key)
}

func (h GinCtxServiceImpl) processBindError(err error) apperrors.BadRequestError {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		var appValErrors []apperrors.ValidationError
		for _, fieldError := range fieldErrors {
			appValError := apperrors.ValidationError{Field: fieldError.Field(), Message: fieldError.ActualTag()}
			appValErrors = append(appValErrors, appValError)
		}
		return apperrors.NewBadReqErrorFromValidationErrors(appValErrors)
	}
	log.WithError(err).Info("Unable to bind json")
	cannotBindJsonRuleErr := h.errorMessageService.RuleErrorFromCode(apperrors.ErrCodeCannotBindJson)
	return apperrors.NewBadReqErrorFromRuleError(cannotBindJsonRuleErr)
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

// NewGinCtxService creates a new GinCtxService instance
func NewGinCtxService(errorService errorservices.ErrorService) GinCtxService {
	return &GinCtxServiceImpl{errorMessageService: errorService}
}

// ReadValueFromBody reads the request body from a gin context and binds it to a
// value provided in the type parameter. Using a pointer type is not permitted
// and will trigger a panic.
func ReadValueFromBody[V any](ginCtxService GinCtxService, c *gin.Context) single.Single[V] {
	var value V
	if err := c.ShouldBind(&value); err != nil {
		return single.Error[V](ginCtxService.processBindError(err))
	}
	return single.Just(value)
}

// BindBodyToPointer reads the request type and writes it to a value referenced by the pointer provided.
func BindBodyToPointer[V any](ginCtxService GinCtxService, c *gin.Context, value *V) single.Single[*V] {
	if err := c.ShouldBind(value); err != nil {
		return single.Error[*V](ginCtxService.processBindError(err))
	}
	return single.Just(value)
}
