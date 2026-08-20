package article

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	group := router.Group("/article")
	group.POST("/", handler.Create)
	group.GET("/*path", handler.DispatchGet)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
