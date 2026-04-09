package agent_server

import (
	"github.com/gin-gonic/gin"
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server/handler"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestLogger(), middleware.Recovery())

	h := handler.NewHandler()
	router.GET("/health", h.Health)

	api := router.Group("/api")
	api.POST("/token", h.Token)

	protected := api.Group("")
	protected.Use(middleware.NeedCowAuth())
	protected.POST("/metric", h.Metric)

	return router
}
