package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/externalservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
)

type UserController interface {
	webservices.Controller
}

type userControllerImpl struct {
	userService    externalservices.ExtUserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  webservices.GinCtxService
}

func (u userControllerImpl) AddRoutes(r *gin.Engine) {
	userGroupV1 := r.Group("test", u.authMiddleware.Authentication())

	userGroupV1.GET("/user",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			identity := security.GetIdentityFromGinContext(c)
			userDto, err := single.RetrieveValue(c, u.userService.GetByAuthId(c, identity.GetAuthId()))
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})
}

func NewUserController(
	authMiddleware middlewares.AuthMiddleware,
	userService externalservices.ExtUserService,
	ginCtxService webservices.GinCtxService,
) UserController {
	return &userControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
	}
}
