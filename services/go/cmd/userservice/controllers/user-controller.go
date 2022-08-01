package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web"
	"github.com/obenkenobi/cypher-log/services/go/pkg/web/webservices"
)

type UserController interface {
	web.Controller
}

type userControllerImpl struct {
	userService    services.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  webservices.GinCtxService
}

func (u userControllerImpl) AddRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")

	userGroup.POST("",
		u.authMiddleware.Authentication(),
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

	userGroup.PUT("",
		u.authMiddleware.Authentication(),
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

	userGroup.DELETE("",
		u.authMiddleware.Authentication(),
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			updateUserSrc := u.userService.DeleteUser(c, security.GetIdentityFromGinContext(c))
			userDto, err := single.RetrieveValue(c, updateUserSrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})

	userGroup.GET("/me",
		u.authMiddleware.Authentication(),
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			getUserIdentitySrc := u.userService.GetUserIdentity(c, security.GetIdentityFromGinContext(c))
			userDto, err := single.RetrieveValue(c, getUserIdentitySrc)
			u.ginCtxService.RespondJsonOkOrError(c, userDto, err)
		})

	userGroup.GET("/byAuthId/:id",
		u.authMiddleware.Authentication(),
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
