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
			func(c *gin.Context) {
				var userCreateDto userdtos.UserSaveDto
				if err := c.ShouldBind(&userCreateDto); err != nil {
					errors.HandleBindError(c, err)
					return
				}
				userDto, err := u.userService.AddUser(userCreateDto)
				if err != nil {
					errors.HandleErrorResponse(c, *err)
					return
				}
				c.JSON(http.StatusOK, userDto)
			})

		userGroup.PUT("",
			func(c *gin.Context) {
				var userCreateDto userdtos.UserSaveDto
				if err := c.ShouldBind(&userCreateDto); err != nil {
					errors.HandleBindError(c, err)
					return
				}
				userDto, err := u.userService.UpdateUser(userCreateDto)
				if err != nil {
					errors.HandleErrorResponse(c, *err)
					return
				}
				c.JSON(http.StatusOK, userDto)
			})

		userGroup.GET("/me",
			u.authMiddleware.Authentication(),
			u.authMiddleware.Authorization(middlewares.AuthorizerMiddlewareParams{RequireAnyAuthorities: []string{"amin"}}),
			func(c *gin.Context) {
				identity := security.NewIdentityHolderFromContext(c)
				log.Info("Identity created", identity)
				userDto := u.userService.GetByProviderUserId(identity.GetSubject())
				c.JSON(http.StatusOK, userDto)
			})

		userGroup.GET("/byProviderUserId/:id",
			func(c *gin.Context) {
				id := c.Param("id")
				userDto := u.userService.GetByProviderUserId(id)
				c.JSON(http.StatusOK, userDto)
			})
	}
}

func NewUserController(authMiddleware middlewares.AuthMiddleware, userService services.UserService) UserController {
	return &userControllerImpl{
		authMiddleware: authMiddleware,
		userService:    userService,
	}
}
