package api

import (
	"github.com/labstack/echo/v4"
)

func config_heart_beat_callback(httpServer *echo.Echo) {
	httpServer.POST("/api/node_info", addPullZone, middleware.ParseToken)
}
