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
			reqUserSrc := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[commondtos.ExistsDto] {
					return u.userKeyService.UserKeyExists(c, userBos)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userKeyGroupV1.POST("/passcode",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[keydtos.PasscodeCreateDto](u.ginCtxService, c)
			if err != nil {
				u.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[commondtos.SuccessDto] {
					return u.userKeyService.CreateUserKey(c, userBos, body)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userKeyGroupV1.POST("/newSession",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[keydtos.PasscodeDto](u.ginCtxService, c)
			if err != nil {
				u.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[commondtos.UKeySessionDto] {
					return u.userKeyService.NewKeySession(c, userBos, body)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOk(c, resBody, err)
		})

	userKeyGroupV1.POST("/getKeyFromSession",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			body, err := ginservices.ReadValueFromBody[commondtos.UKeySessionDto](u.ginCtxService, c)
			if err != nil {
				u.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			resBody, err := single.RetrieveValue(c, u.userKeyService.GetKeyFromSession(c, body))
			u.ginCtxService.RespondJsonOk(c, resBody, err)
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
