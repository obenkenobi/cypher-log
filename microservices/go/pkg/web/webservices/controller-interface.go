package webservices

import "github.com/gin-gonic/gin"

// Controller is an interface that contains baseline methods needed by all
// controllers in a web application. (Recall the Model-View-Controller pattern.)
type Controller interface {
	// AddRoutes adds http routes to a Gin Engine instance.
	AddRoutes(r *gin.Engine)
}
