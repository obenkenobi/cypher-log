package ginservices

import (
	"github.com/akrennmair/slice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/queryreq"
	"net/http"
)

// GinCtxService is an injectable service that provides convenience methods
// relating to http requests and responses using the Gin framework. As the name
// suggest, it is designed to handle a Gin's Context struct.
type GinCtxService interface {
	// ReqQueryReader returns a reader to parse your request query params
	ReqQueryReader(c *gin.Context) queryreq.ReqQueryReader

	// RespondError takes an error and parses it to set the appropriate http
	// response. Certain errors relating to user input will trigger a 4XX status
	// code. Otherwise, a 5XX code will be thrown indicating the error means
	// something went wrong with the server.
	RespondError(c *gin.Context, err error)

	// processBindError takes an error from binding a value from a request body
	// processes it into a BadRequestError.
	processBindError(err error) apperrors.BadRequestError

	// RestControllerPipeline initializes a controller.Pipeline that manages errors in the
	// background using the gin context for REST controllers.
	RestControllerPipeline(c *gin.Context) controller.Pipeline
}

type GinCtxServiceImpl struct {
	errorMessageService sharedservices.ErrorService
}

func (g GinCtxServiceImpl) ReqQueryReader(c *gin.Context) queryreq.ReqQueryReader {
	return queryreq.NewGinCtxReqQueryReaderImpl(c, g.errorMessageService)
}

func (g GinCtxServiceImpl) RestControllerPipeline(c *gin.Context) controller.Pipeline {
	return controller.NewPipelineImpl(func(err error) {
		g.RespondError(c, err)
	})
}

func (g GinCtxServiceImpl) RespondError(c *gin.Context, err error) {
	if badReqErr, ok := err.(apperrors.BadRequestError); ok {
		c.JSON(http.StatusBadRequest, badReqErr)
	} else {
		logger.Log.Error(err)
		c.Status(http.StatusInternalServerError)
	}
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

// BindBodyToReferenceObj reads the request type and writes it to a value that must be a reference type
func BindBodyToReferenceObj[V any](ginCtxService GinCtxService, c *gin.Context, value V) (V, error) {
	if err := c.ShouldBind(value); err != nil {
		return value, ginCtxService.processBindError(err)
	}
	return value, nil
}
