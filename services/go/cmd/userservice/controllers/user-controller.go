package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/ginextensions"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"net/http"
)

type UserController interface {
	AddRoutes(r *gin.Engine)
}

type userControllerImpl struct {
	userService       services.UserService
	authMiddleware    middlewares.AuthMiddleware
	ginWrapperService ginextensions.GinWrapperService
}

func (u userControllerImpl) AddRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				bindBodyX := ginextensions.BindValueToBody(u.ginWrapperService, c, userdtos.UserSaveDto{})
				addUserX := stream.FlatMap(bindBodyX,
					func(userSaveDto userdtos.UserSaveDto) stream.Observable[userdtos.UserDto] {
						return u.userService.AddUser(security.GetIdentityFromContext(c), userSaveDto)
					})
				if userDto, err := stream.First(c, addUserX); err != nil {
					u.ginWrapperService.HandleErrorResponse(c, err)
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})

		userGroup.PUT("",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				bindBodyX := ginextensions.BindValueToBody(u.ginWrapperService, c, userdtos.UserSaveDto{})
				updateUserX := stream.FlatMap(bindBodyX,
					func(userSaveDto userdtos.UserSaveDto) stream.Observable[userdtos.UserDto] {
						return u.userService.UpdateUser(security.GetIdentityFromContext(c), userSaveDto)
					})
				userDto, err := stream.First(c, updateUserX)
				u.ginWrapperService.RespondJsonOk(c, userDto, err)
			})

		userGroup.GET("/me",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				getUserIdentityX := u.userService.GetUserIdentity(security.GetIdentityFromContext(c))
				userDto, err := stream.First(c, getUserIdentityX)
				u.ginWrapperService.RespondJsonOk(c, userDto, err)
			})

		userGroup.GET("/byAuthId/:id",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
			func(c *gin.Context) {
				authId := c.Param("id")
				getUserX := u.userService.GetByAuthId(authId)
				userDto, err := stream.First(c, getUserX)
				u.ginWrapperService.RespondJsonOk(c, userDto, err)
			})
	}
}

func NewUserController(
	authMiddleware middlewares.AuthMiddleware,
	userService services.UserService,
	ginWrapperService ginextensions.GinWrapperService,
) UserController {
	return &userControllerImpl{
		authMiddleware:    authMiddleware,
		userService:       userService,
		ginWrapperService: ginWrapperService,
	}
}
