package web_server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/web_server/handler"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(corsMiddleware(), middleware.RequestLogger(), middleware.Recovery())

	h := handler.NewHandler()
	router.GET("/health", h.Health)

	api := router.Group("/api")
	registerAPIRoutes(api, h)

	return router
}

func registerAPIRoutes(api *gin.RouterGroup, h *handler.Handler) {
	registerPublicRoutes(api, h)

	protected := api.Group("")
	protected.Use(middleware.NeedAuth())
	registerProtectedRoutes(protected, h)
}

func registerPublicRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/auth")
	group.POST("/login", h.AuthLogin)
	group.POST("/password-reset/request", h.AuthPasswordResetRequest)
	group.POST("/password-reset/verify", h.AuthPasswordResetVerify)
	group.POST("/password-reset/confirm", h.AuthPasswordResetConfirm)
}

func registerProtectedRoutes(api *gin.RouterGroup, h *handler.Handler) {
	registerAuthRoutes(api, h)
	registerDashboardRoutes(api, h)
	registerCowRoutes(api, h)
	registerUserRoutes(api, h)
	registerReportRoutes(api, h)
	registerAlertRoutes(api, h)
	registerMetricRoutes(api, h)
}

func registerAuthRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/auth")
	group.POST("/refresh", h.AuthRefreshToken)
}

func registerDashboardRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/dashboard")
	group.GET("/summary", h.DashboardSummary)
	group.GET("/list", h.DashboardList)
}

func registerCowRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/cow")
	group.GET("/list", h.CowList)
	group.GET("/info", h.CowInfo)
	group.POST("/create", h.CowCreate)
	group.POST("/update", h.CowUpdate)

	metricGroup := group.Group("/metric")
	metricGroup.GET("/temperature", h.CowMetricTemperature)
	metricGroup.GET("/heart_rate", h.CowMetricHeartRate)
	metricGroup.GET("/blood_oxygen", h.CowMetricBloodOxygen)
	metricGroup.GET("/milk_amount", h.CowMetricMilkAmount)
	metricGroup.GET("/movement", h.CowMetricMovement)
	metricGroup.GET("/movement_path", h.CowMetricMovementPath)
	metricGroup.GET("/weight", h.CowMetricWeight)
}

func registerUserRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/user")
	group.GET("/list", h.UserList)
	group.GET("/info", h.UserInfo)
	group.POST("/create", h.UserCreate)
	group.POST("/update", h.UserUpdate)
	group.POST("/update_password", h.UserUpdatePassword)
	group.POST("/delete", h.UserDelete)
}

func registerReportRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/report")
	group.GET("/list", h.ReportList)
	group.GET("/latest", h.ReportLatest)
}

func registerAlertRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/alert")
	group.GET("/summary", h.AlertSummary)
	group.GET("/list", h.AlertList)
}

func registerMetricRoutes(api *gin.RouterGroup, h *handler.Handler) {
	group := api.Group("/metric")
	group.GET("/list", h.MetricList)
}
