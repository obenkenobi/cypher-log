package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/middlewares"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	cDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	nDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/notedtos"
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

	noteGroupV1.POST("",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto]](n.ginCtxService, c)
			if err != nil {
				n.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc, func(userBos userbos.UserBo) single.Single[cDTOs.SuccessDto] {
				return n.noteService.AddNoteTransaction(c, userBos, body)
			},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			n.ginCtxService.RespondJsonOk(c, resBody, err)
		})
	noteGroupV1.PUT("",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto]](n.ginCtxService, c)
			if err != nil {
				n.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[cDTOs.SuccessDto] {
					return n.noteService.UpdateNoteTransaction(c, userBos, body)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			n.ginCtxService.RespondJsonOk(c, resBody, err)
		})
	noteGroupV1.DELETE("",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[nDTOs.NoteIdDto](n.ginCtxService, c)
			if err != nil {
				n.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[cDTOs.SuccessDto] {
					return n.noteService.DeleteNoteTransaction(c, userBos, body)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			n.ginCtxService.RespondJsonOk(c, resBody, err)
		})
	noteGroupV1.POST("/getById",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto]](n.ginCtxService, c)
			if err != nil {
				n.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[nDTOs.NoteReadDto] {
					return n.noteService.GetNoteById(c, userBos, body)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			n.ginCtxService.RespondJsonOk(c, resBody, err)
		})
	noteGroupV1.POST("/getPage",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			reqUserSrc := n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
			body, err := ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[pagination.PageRequest]](n.ginCtxService, c)
			if err != nil {
				n.ginCtxService.HandleErrorResponse(c, err)
				return
			}
			businessLogicSrc := single.FlatMap(reqUserSrc,
				func(userBos userbos.UserBo) single.Single[pagination.Page[nDTOs.NotePreviewDto]] {
					return n.noteService.GetNotesPage(c, userBos, body)
				},
			)
			resBody, err := single.RetrieveValue(c, businessLogicSrc)
			n.ginCtxService.RespondJsonOk(c, resBody, err)
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
