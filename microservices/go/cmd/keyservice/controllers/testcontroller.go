package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
)

type TestController interface {
	webservices.Controller
}

type testControllerImpl struct {
	userService    services.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  webservices.GinCtxService
}

func (u testControllerImpl) AddRoutes(r *gin.Engine) {
	userGroupV1 := r.Group("test", u.authMiddleware.Authentication())

	userGroupV1.GET("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUser := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c)).ScheduleAsync(c)
			user, err := single.RetrieveValue(c, reqUser)
			u.ginCtxService.RespondJsonOkOrError(c, user, err)
		})
}

func NewTestController(
	authMiddleware middlewares.AuthMiddleware,
	userService services.UserService,
	ginCtxService webservices.GinCtxService,
) TestController {
	return &testControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
	}
}
