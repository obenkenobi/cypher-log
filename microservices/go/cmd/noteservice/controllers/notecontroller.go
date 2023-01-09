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
	"net/http"
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
			var userBo userbos.UserBo
			var reqBody cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto]
			var resBody cDTOs.SuccessDto

			n.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[nDTOs.NoteCreateDto]](
					n.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, n.noteService.AddNoteTransaction(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})
	noteGroupV1.PUT("",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var reqBody cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto]
			var resBody cDTOs.SuccessDto

			n.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[nDTOs.NoteUpdateDto]](
					n.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, n.noteService.UpdateNoteTransaction(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})
	noteGroupV1.DELETE("",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var reqBody nDTOs.NoteIdDto
			var resBody cDTOs.SuccessDto

			n.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[nDTOs.NoteIdDto](n.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, n.noteService.DeleteNoteTransaction(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})
	noteGroupV1.POST("/getById",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var reqBody cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto]
			var resBody nDTOs.NoteReadDto

			n.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[nDTOs.NoteIdDto]](
					n.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, n.noteService.GetNoteById(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
		})
	noteGroupV1.POST("/getPage",
		n.authMiddleware.Authorization(middlewares.AuthorizerSettings{VerifyIsUser: true}),
		func(c *gin.Context) {
			var userBo userbos.UserBo
			var reqBody cDTOs.UKeySessionReqDto[pagination.PageRequest]
			var resBody pagination.Page[nDTOs.NotePreviewDto]

			n.ginCtxService.StartCtxPipeline(c).Next(func() (err error) {
				userBo, err = n.userService.RequireUser(c, security.GetIdentityFromGinContext(c))
				return
			}).Next(func() (err error) {
				reqBody, err = ginservices.ReadValueFromBody[cDTOs.UKeySessionReqDto[pagination.PageRequest]](
					n.ginCtxService, c)
				return
			}).Next(func() (err error) {
				resBody, err = single.RetrieveValue(c, n.noteService.GetNotesPage(c, userBo, reqBody))
				return
			}).Next(func() (err error) {
				c.JSON(http.StatusOK, resBody)
				return
			})
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
