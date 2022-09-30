package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
)

type TestController interface {
	controller.Controller
}

type testControllerImpl struct {
	userService    sharedservices.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  ginservices.GinCtxService
}

func (u testControllerImpl) AddRoutes(r *gin.Engine) {
	userGroupV1 := r.Group("test", u.authMiddleware.Authentication())

	userGroupV1.GET("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUser := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c)).ScheduleEagerAsync(c)
			user, err := single.RetrieveValue(c, reqUser)
			u.ginCtxService.RespondJsonOkOrError(c, user, err)
		})
}

func NewTestController(
	authMiddleware middlewares.AuthMiddleware,
	userService sharedservices.UserService,
	ginCtxService ginservices.GinCtxService,
) TestController {
	return &testControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
	}
}
