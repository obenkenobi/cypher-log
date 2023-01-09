package ginservices

import (
	"github.com/akrennmair/slice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/queryreq"
	"net/http"
)

// GinCtxService is an injectable service that provides convenience methods
// relating to http requests and responses using the Gin framework. As the name
// suggest, it is designed to handle a Gin's Context struct.
type GinCtxService interface {
	// ParamStr returns the value of a URL param as a string
	ParamStr(c *gin.Context, name string) string

	// ReqQueryReader returns a reader to parse your request query params
	ReqQueryReader(c *gin.Context) queryreq.ReqQueryReader

	// HandleErrorResponse takes an error and parses it to set the appropriate http
	// response. Certain errors relating to user input will trigger a 4XX status
	// code. Otherwise, a 5XX code will be thrown indicating the error means
	// something went wrong with the server.
	HandleErrorResponse(c *gin.Context, err error)

	// RespondJsonOk responds with json value with a 200 status code or an
	// error if the error != nil.
	RespondJsonOk(c *gin.Context, model any, err error)

	// processBindError takes an error from binding a value from a request body processes it into a BadRequestError.
	processBindError(err error) apperrors.BadRequestError
}

type GinCtxServiceImpl struct {
	errorMessageService sharedservices.ErrorService
}

func (g GinCtxServiceImpl) ParamStr(c *gin.Context, key string) string {
	return c.Param(key)
}

func (q GinCtxServiceImpl) ReqQueryReader(c *gin.Context) queryreq.ReqQueryReader {
	return queryreq.NewGinCtxReqQueryReaderImpl(c, q.errorMessageService)
}

func (g GinCtxServiceImpl) HandleErrorResponse(c *gin.Context, err error) {
	if badReqErr, ok := err.(apperrors.BadRequestError); ok {
		c.JSON(http.StatusBadRequest, badReqErr)
	} else {
		logger.Log.Error(err)
		c.Status(http.StatusInternalServerError)
	}
}

func (g GinCtxServiceImpl) RespondJsonOk(c *gin.Context, model any, err error) {
	if err != nil {
		g.HandleErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, model)
}

func (g GinCtxServiceImpl) processBindError(err error) apperrors.BadRequestError {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		appValErrors := slice.Map(fieldErrors, func(fieldError validator.FieldError) apperrors.ValidationError {
			return apperrors.ValidationError{Field: fieldError.Field(), Message: fieldError.ActualTag()}
		})
		return apperrors.NewBadReqErrorFromValidationErrors(appValErrors)
	}
	logger.Log.WithError(err).Info("Unable to bind json")
	cannotBindJsonRuleErr := g.errorMessageService.RuleErrorFromCode(apperrors.ErrCodeCannotBindJson)
	return apperrors.NewBadReqErrorFromRuleError(cannotBindJsonRuleErr)
}

// NewGinCtxServiceImpl creates a new GinCtxService instance
func NewGinCtxServiceImpl(errorService sharedservices.ErrorService) *GinCtxServiceImpl {
	return &GinCtxServiceImpl{errorMessageService: errorService}
}

// ReadValueFromBody reads the request body from a gin context and binds it to a
// value provided in the type parameter. Using a pointer type is not permitted
// and will trigger a panic.
func ReadValueFromBody[V any](ginCtxService GinCtxService, c *gin.Context) (V, error) {
	var value V
	if err := c.ShouldBind(&value); err != nil {
		return value, ginCtxService.processBindError(err)
	}
	return value, nil
}

// BindBodyToPointer reads the request type and writes it to a value referenced by the pointer provided.
func BindBodyToPointer[V any](ginCtxService GinCtxService, c *gin.Context, value *V) (*V, error) {
	if err := c.ShouldBind(value); err != nil {
		return value, ginCtxService.processBindError(err)
	}
	return value, nil
}
