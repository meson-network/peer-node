package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/node_info"
	"github.com/meson-network/peer_common/heart_beat"
)

func config_heart_beat_callback(httpServer *echo_plugin.EchoServer) {
	httpServer.GET("/api/node_info", nodeInfoHandler, MidToken)
}

// @Summary      /api/health
// @Description  health check
// @Tags         health
// @Produce      json
// @Success      200 {object} MSG_RESP_HEALTH "server unix time"
// @Router       /api/node_info [get]
func nodeInfoHandler(ctx echo.Context) error {
	res := &heart_beat.Msg_Resp_HeartBeatCallback{}

	//check token
	err := CheckToken(ctx)
	if err != nil {
		res.MetaStatus(-2, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	nodeInfo := node_info.GetNodeInfo()
	res.Stor_total_bytes = nodeInfo.Stor_total_bytes
	res.Stor_used_bytes = nodeInfo.Stor_used_bytes
	res.HardwareInfo = nodeInfo.HardwareInfo
	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}
