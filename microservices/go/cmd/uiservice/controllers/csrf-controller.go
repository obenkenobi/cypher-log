package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/web/controller"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

type CsrfController interface {
	controller.Controller
}

type CsrfControllerImpl struct {
}

func (c CsrfControllerImpl) AddRoutes(r *gin.Engine) {
	r.GET("/csrf", func(c *gin.Context) {
		c.SetCookie("XSRF-TOKEN", csrf.GetToken(c), 0, "", "", true, false)
		c.JSON(http.StatusOK, commondtos.NewSuccessTrue())
	})
}

func NewCsrfControllerImpl() *CsrfControllerImpl {
	return &CsrfControllerImpl{}
}
