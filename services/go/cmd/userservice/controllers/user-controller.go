package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/ginx"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/streamx/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"github.com/obenkenobi/cypher-log/services/go/pkg/server"
	"net/http"
)

type UserController interface {
	server.Controller
}

type userControllerImpl struct {
	userService    services.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  ginx.GinCtxService
}

func (u userControllerImpl) AddRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")

	userGroup.POST("",
		u.authMiddleware.Authentication(),
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			bindBodySrc := ginx.BindValueToBody(u.ginCtxService, c, userdtos.UserSaveDto{})
			addUserSrc := single.FlatMap(bindBodySrc,
				func(userSaveDto userdtos.UserSaveDto) single.Single[userdtos.UserDto] {
					return u.userService.AddUser(security.GetIdentityFromContext(c), userSaveDto)
				})
			if userDto, err := single.AwaitItem(c, addUserSrc); err != nil {
				u.ginCtxService.HandleErrorResponse(c, err)
			} else {
				c.JSON(http.StatusOK, userDto)
			}
		})

	userGroup.PUT("",
		u.authMiddleware.Authentication(),
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			bindBodySrc := ginx.BindValueToBody(u.ginCtxService, c, userdtos.UserSaveDto{})
			updateUserSrc := single.FlatMap(bindBodySrc,
				func(userSaveDto userdtos.UserSaveDto) single.Single[userdtos.UserDto] {
					return u.userService.UpdateUser(security.GetIdentityFromContext(c), userSaveDto)
				})
			userDto, err := single.AwaitItem(c, updateUserSrc)
			u.ginCtxService.RespondJsonOk(c, userDto, err)
		})

	userGroup.GET("/me",
		u.authMiddleware.Authentication(),
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			getUserIdentitySrc := u.userService.GetUserIdentity(security.GetIdentityFromContext(c))
			userDto, err := single.AwaitItem(c, getUserIdentitySrc)
			u.ginCtxService.RespondJsonOk(c, userDto, err)
		})

	userGroup.GET("/byAuthId/:id",
		u.authMiddleware.Authentication(),
		u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
		func(c *gin.Context) {
			authId := c.Param("id")
			getUserSrc := u.userService.GetByAuthId(authId)
			userDto, err := single.AwaitItem(c, getUserSrc)
			u.ginCtxService.RespondJsonOk(c, userDto, err)
		})

}

func NewUserController(
	authMiddleware middlewares.AuthMiddleware,
	userService services.UserService,
	ginCtxService ginx.GinCtxService,
) UserController {
	return &userControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
		ginCtxService:  ginCtxService,
	}
}
