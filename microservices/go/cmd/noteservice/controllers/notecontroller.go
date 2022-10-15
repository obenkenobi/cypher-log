package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/security"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/ginservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/routing"
)

type NoteController interface {
	controller.Controller
}

type NoteControllerImpl struct {
	userService    sharedservices.UserService
	authMiddleware middlewares.AuthMiddleware
	ginCtxService  ginservices.GinCtxService
	noteService    services.NoteService
}

func (n NoteControllerImpl) AddRoutes(r *gin.Engine) {
	noteGroupV1 := r.Group(routing.APIPath(1, "notes"), n.authMiddleware.Authentication())

	noteGroupV1.GET("/test",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			resBody, err := single.RetrieveValue(c, reqUserSrc)
			n.ginCtxService.RespondJsonOkOrError(c, resBody, err)
		})
}

func NewNoteControllerImpl(
	userService sharedservices.UserService,
	authMiddleware middlewares.AuthMiddleware,
	ginCtxService ginservices.GinCtxService,
	noteService services.NoteService,
) *NoteControllerImpl {
	return &NoteControllerImpl{
		userService:    userService,
		authMiddleware: authMiddleware,
		ginCtxService:  ginCtxService,
		noteService:    noteService,
	}
}
