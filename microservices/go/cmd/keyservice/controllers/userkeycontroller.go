package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/routing"
	"net/http"
)

type UserKeyController interface {
	controller.Controller
}

type UserKeyControllerImpl struct {
	userService    sharedservices.UserService
	userKeyService services.UserKeyService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  ginservices.GinCtxService
}

func (u UserKeyControllerImpl) AddRoutes(r *gin.Engine) {
	userKeyGroupV1 := r.Group(routing.APIPath(1, "userKey"), u.authMiddleware.Authentication())

	userKeyGroupV1.GET("/exists",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var resBody commondtos.ExistsDto

			u.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, u.userKeyService.UserKeyExists(c, userBo))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})

	userKeyGroupV1.POST("/passcode",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var reqBody keydtos.PasscodeCreateDto
			var resBody commondtos.SuccessDto

			u.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[keydtos.PasscodeCreateDto](u.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, u.userKeyService.CreateUserKey(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})

	userKeyGroupV1.POST("/newSession",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var reqBody keydtos.PasscodeDto
			var resBody commondtos.UKeySessionDto

			u.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[keydtos.PasscodeDto](u.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, u.userKeyService.NewKeySession(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})

	userKeyGroupV1.POST("/getKeyFromSession",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			var reqBody commondtos.UKeySessionDto
			var resBody keydtos.UserKeyDto

			u.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[commondtos.UKeySessionDto](u.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, u.userKeyService.GetKeyFromSession(c, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})
}

func NewTestControllerImpl(
	authMiddleware middlewares.AuthMiddleware,
	userService sharedservices.UserService,
	ginCtxService ginservices.GinCtxService,
	userKeyService services.UserKeyService,
) *UserKeyControllerImpl {
	return &UserKeyControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
		userKeyService: userKeyService,
	}
}
