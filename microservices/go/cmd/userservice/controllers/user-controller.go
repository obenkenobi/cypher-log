package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/routing"
	"net/http"
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
			var reqBody userdtos.UserSaveDto
			var resBody userdtos.UserReadDto

			u.ginCtxService.RestControllerPipeline(c).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[userdtos.UserSaveDto](u.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = u.userService.AddUserTxn(c, security.GetIdentityFromGinContext(c), reqBody)
				return
			}).Next(func() error {
				c.JSON(http.StatusOK, resBody)
				return nil
			})
		})

	userGroupV1.PUT("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var reqBody userdtos.UserSaveDto
			var resBody userdtos.UserReadDto

			u.ginCtxService.RestControllerPipeline(c).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[userdtos.UserSaveDto](u.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = u.userService.UpdateUserTxn(c, security.GetIdentityFromGinContext(c), reqBody)
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return nil
			})
		})

	userGroupV1.DELETE("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var resBody userdtos.UserReadDto

			u.ginCtxService.RestControllerPipeline(c).Next(func() (err error) {
				resBody, err = u.userService.BeginDeletingUserTxn(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})

	userGroupV1.GET("/me",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var resBody userdtos.UserIdentityDto

			u.ginCtxService.RestControllerPipeline(c).Next(func() (err error) {
				resBody, err = u.userService.GetUserIdentity(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})

	userGroupV1.GET("/byAuthId/:id",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			var resBody userdtos.UserReadDto
			authId := c.Param("id")

			u.ginCtxService.RestControllerPipeline(c).Next(func() (err error) {
				resBody, err = u.userService.GetByAuthId(c, authId)
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})

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
