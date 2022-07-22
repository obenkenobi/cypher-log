package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/services/go/pkg/security"
	log "github.com/sirupsen/logrus"
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
				userSaveDto := userdtos.UserSaveDto{}
				if err := c.ShouldBind(&userSaveDto); err != nil {
					errors.HandleBindError(c, err)
					return
				}
				if userDto, err := u.userService.AddUser(userSaveDto); err != nil {
					errors.HandleErrorResponse(c, *err)
					return
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})

		userGroup.PUT("",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				userSaveDto := userdtos.UserSaveDto{}
				if err := c.ShouldBind(&userSaveDto); err != nil {
					errors.HandleBindError(c, err)
					return
				}
				identityHolder := security.NewIdentityHolderFromContext(c)
				if userDto, err := u.userService.UpdateUser(identityHolder, userSaveDto); err != nil {
					errors.HandleErrorResponse(c, *err)
					return
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})

		userGroup.GET("/me",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
			func(c *gin.Context) {
				identity := security.NewIdentityHolderFromContext(c)
				log.Info("Identity created", identity)
				if userDto, err := u.userService.GetByProviderUserId(identity.GetIdFromProvider()); err != nil {
					errors.HandleErrorResponse(c, *err)
					return
				} else {
					c.JSON(http.StatusOK, userDto)
				}
			})

		userGroup.GET("/byProviderUserId/:id",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsSystemClient: true}),
			func(c *gin.Context) {
				id := c.Param("id")
				if userDto, err := u.userService.GetByProviderUserId(id); err != nil {
					errors.HandleErrorResponse(c, *err)
					return
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
