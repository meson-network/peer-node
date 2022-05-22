package api

import (
	"github.com/labstack/echo/v4"
)

func config_file(httpServer *echo.Echo) {
	httpServer.POST("/api/file/download_task", addPullZone, middleware.ParseToken)
	httpServer.POST("/*", queryPullZone, middleware.ParseToken)
}
