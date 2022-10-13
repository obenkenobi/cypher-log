package controllers

import (
	"github.com/barweiss/go-tuple"
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
			u.ginCtxService.RespondJsonOkOrError(c, resBody, err)
		})

	userKeyGroupV1.POST("/passcode",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			bodySrc := ginservices.ReadValueFromBody[keydtos.PasscodeCreateDto](u.ginCtxService, c)
			businessLogicSrc := single.FlatMap(single.Zip2(reqUserSrc, bodySrc),
				func(t tuple.T2[userbos.UserBo, keydtos.PasscodeCreateDto]) single.Single[commondtos.SuccessDto] {
					userBos, passcodeDto := t.V1, t.V2
					return u.userKeyService.CreateUserKey(c, userBos, passcodeDto)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOkOrError(c, resBody, err)
		})

	userKeyGroupV1.POST("/newSession",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := u.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			bodySrc := ginservices.ReadValueFromBody[keydtos.PasscodeDto](u.ginCtxService, c)
			businessLogicSrc := single.FlatMap(single.Zip2(reqUserSrc, bodySrc),
				func(t tuple.T2[userbos.UserBo, keydtos.PasscodeDto]) single.Single[keydtos.UserKeySessionDto] {
					userBos, passcodeDto := t.V1, t.V2
					return u.userKeyService.NewKeySession(c, userBos, passcodeDto)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOkOrError(c, resBody, err)
		})

	userKeyGroupV1.POST("/getKeyFromSession",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			bodySrc := ginservices.ReadValueFromBody[keydtos.UserKeySessionDto](u.ginCtxService, c)
			businessLogicSrc := single.FlatMap(bodySrc,
				func(sessionDto keydtos.UserKeySessionDto) single.Single[keydtos.UserKeyDto] {
					return u.userKeyService.GetKeyFromSession(c, sessionDto)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			u.ginCtxService.RespondJsonOkOrError(c, resBody, err)
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
