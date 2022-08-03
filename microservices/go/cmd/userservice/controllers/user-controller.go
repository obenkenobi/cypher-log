package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/webservices"
)

type UserController interface {
	webservices.Controller
}

type userControllerImpl struct {
	userService    services.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  webservices.GinCtxService
}

func (u userControllerImpl) AddRoutes(r *gin.Engine) {
	userGroupV1 := r.Group(web.APIPath(1, "user"), u.authMiddleware.Authentication())

	userGroupV1.POST("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			readValFromBodySrc := webservices.ReadValueFromBody[userdtos.UserSaveDto](u.ginCtxService, c)
			addUserSrc := single.FlatMap(readValFromBodySrc,
				func(userSaveDto userdtos.UserSaveDto) single.Single[userdtos.UserDto] {
					return u.userService.AddUser(c, security.GetIdentityFromGinContext(c), userSaveDto)
				})
			userDto, err := single.RetrieveValue(c, addUserSrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})

	userGroupV1.PUT("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			readValFromBodySrc := webservices.ReadValueFromBody[userdtos.UserSaveDto](u.ginCtxService, c)
			updateUserSrc := single.FlatMap(readValFromBodySrc,
				func(userSaveDto userdtos.UserSaveDto) single.Single[userdtos.UserDto] {
					return u.userService.UpdateUser(c, security.GetIdentityFromGinContext(c), userSaveDto)
				})
			userDto, err := single.RetrieveValue(c, updateUserSrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})

	userGroupV1.DELETE("",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			updateUserSrc := u.userService.DeleteUser(c, security.GetIdentityFromGinContext(c))
			userDto, err := single.RetrieveValue(c, updateUserSrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})

	userGroupV1.GET("/me",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			getUserIdentitySrc := u.userService.GetUserIdentity(c, security.GetIdentityFromGinContext(c))
			userDto, err := single.RetrieveValue(c, getUserIdentitySrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})

	userGroupV1.GET("/byAuthId/:id",
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			authId := c.Param("id")
			getUserSrc := u.userService.GetByAuthId(c, authId)
			userDto, err := single.RetrieveValue(c, getUserSrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})
}

func NewUserController(
	authMiddleware middlewares.AuthMiddleware,
	userService services.UserService,
	ginCtxService webservices.GinCtxService,
) UserController {
	return &userControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
	}
}
