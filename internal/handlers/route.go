package handlers

import "github.com/gin-gonic/gin"

type IRouteHandler interface {
	SetUpRoutes(r *gin.RouterGroup)
}
