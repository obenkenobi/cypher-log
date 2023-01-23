package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"io/ioutil"
)

type UserKeyMiddleware interface {
	// UserKeySession wraps restful request bodies with your user-key session if a
	// query param is set to be passUserKeySession=true
	UserKeySession() gin.HandlerFunc
}

type UserKeyMiddlewareImpl struct {
	ginCtxService ginservices.GinCtxService
}

func (u UserKeyMiddlewareImpl) UserKeySession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read query param to determine if the user key session will be passed
		passUserKeySession := false

		err := u.ginCtxService.ReqQueryReader(c).
			ReadBoolOrDefault("passUserKeySession", &passUserKeySession, false).
			Complete()

		if err != nil {
			u.ginCtxService.RespondError(c, err)
			return
		}

		// if user key session will not be passed, this middleware ends
		if !passUserKeySession {
			c.Next()
			return
		}

		// Read the original body if it exists
		originalBody := make(map[string]any)
		if c.Request.Body != nil {
			_, err := ginservices.BindBodyToReferenceObj(u.ginCtxService, c, originalBody)
			if err != nil {
				u.ginCtxService.RespondError(c, err)
				return
			}
		}

		// Get the session dto from your session.
		// If it cannot be read, it means the session is empty and will be used anyway.
		sessionDto, _ := security.ReadUKeySessionDtoFromSession(c)

		// Create a new body as a session request
		newBody := commondtos.UKeySessionReqDto[map[string]any]{
			Session: sessionDto,
			Value:   originalBody,
		}

		// Convert the body into JSON bytes
		newBodyBytes, err := json.Marshal(newBody)
		if err != nil {
			u.ginCtxService.RespondError(c, err)
			return
		}

		// Set the new body onto the request
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(newBodyBytes))
		c.Next()
	}
}

func NewUserKeyMiddlewareImpl(ginCtxService ginservices.GinCtxService) *UserKeyMiddlewareImpl {
	return &UserKeyMiddlewareImpl{ginCtxService: ginCtxService}
}
