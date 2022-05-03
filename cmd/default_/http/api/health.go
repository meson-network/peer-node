package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
)

type MSG_RESP_HEALTH struct {
	UnixTime int64 `json:"unixtime"`
}

// @Summary      /api/health
// @Description  health check
// @Tags         health
// @Produce      json
// @Success      200 {object} MSG_RESP_HEALTH "server unix time"
// @Router       /api/health [get]
func healthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &MSG_RESP_HEALTH{UnixTime: time.Now().Unix()})
}

func config_health(httpServer *echo_plugin.EchoServer) {
	//health
	httpServer.GET("/api/health", healthCheck)
}
