package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/routing"
)

type UserController interface {
	controller.Controller
}

type UserControllerImpl struct {
	userService    services.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  ginservices.GinCtxService
}

func (u UserControllerImpl) AddRoutes(r *gin.Engine) {
	userGroupV1 := r.Group(routing.APIPath(1, "user"), u.authMiddleware.Authentication())

	userGroupV1.POST("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			userSaveDto, err := ginservices.ReadValueFromBody[userdtos.UserSaveDto](u.ginCtxService, c)
			if err != nil {
				u.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			resBody, err := single.RetrieveValue(c,
				u.userService.AddUserTransaction(c, security.GetIdentityFromGinContext(c), userSaveDto))
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userGroupV1.PUT("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			userSaveDto, err := ginservices.ReadValueFromBody[userdtos.UserSaveDto](u.ginCtxService, c)
			if err != nil {
				u.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			resBody, err := single.RetrieveValue(c,
				u.userService.UpdateUserTransaction(c, security.GetIdentityFromGinContext(c), userSaveDto))
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userGroupV1.DELETE("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			businessLogicSrc := u.userService.BeginDeletingUserTransaction(c, security.GetIdentityFromGinContext(c))
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userGroupV1.GET("/me",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			businessLogicSrc := u.userService.GetUserIdentity(c, security.GetIdentityFromGinContext(c))
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userGroupV1.GET("/byAuthId/:id",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			authId := c.Param("id")
			businessLogicSrc := u.userService.GetByAuthId(c, authId)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})
}

func NewUserControllerImpl(
	authMiddleware middlewares.AuthMiddleware,
	userService services.UserService,
	ginCtxService ginservices.GinCtxService,
) *UserControllerImpl {
	return &UserControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
	}
}
