package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	"net/http"
)

type UserController interface {
	AddRoutes(r *gin.Engine)
}

type userControllerImpl struct {
	userService    services.UserService
	authMiddleware middlewares.AuthMiddleware
}

func (u userControllerImpl) AddRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				userSaveDto := &userdtos.UserSaveDto{}
				if err := c.ShouldBind(userSaveDto); err != nil {
					apperrors.HandleBindError(c, err)
					return
				}
				identity := security.GetIdentityFromContext(c)
				if userDto, err := u.userService.AddUser(identity, userSaveDto); err != nil {
					apperrors.HandleErrorResponse(c, *err)
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})

		userGroup.PUT("",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				userSaveDto := &userdtos.UserSaveDto{}
				if err := c.ShouldBind(userSaveDto); err != nil {
					apperrors.HandleBindError(c, err)
					return
				}
				identity := security.GetIdentityFromContext(c)
				if userDto, err := u.userService.UpdateUser(identity, userSaveDto); err != nil {
					apperrors.HandleErrorResponse(c, *err)
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})

		userGroup.GET("/me",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				identity := security.GetIdentityFromContext(c)
				if dto, err := u.userService.GetUserIdentity(identity); err != nil {
					apperrors.HandleErrorResponse(c, *err)
				} else {
					c.JSON(http.StatusOK, dto)
				}
			})

		userGroup.GET("/byProviderUserId/:id",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
			func(c *gin.Context) {
				id := c.Param("id")
				if userDto, err := u.userService.GetByAuthId(id); err != nil {
					apperrors.HandleErrorResponse(c, *err)
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})
	}
}

func NewUserController(authMiddleware middlewares.AuthMiddleware, userService services.UserService) UserController {
	return &userControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
	}
}
