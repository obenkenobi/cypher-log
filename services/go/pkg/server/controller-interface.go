package server

import "github.com/gin-gonic/gin"

type Controller interface {
	AddRoutes(r *gin.Engine)
}
