package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/download_mgr"
	"github.com/meson-network/peer_common/download"
)

func config_download(httpServer *echo_plugin.EchoServer) {
	httpServer.POST("/api/file/download_task", downloadTaskHandler, MidToken)
	//httpServer.POST("/*", queryPullZone, middleware.ParseToken)
}

// @Summary      /api/health
// @Description  health check
// @Tags         health
// @Produce      json
// @Success      200 {object} MSG_RESP_HEALTH "server unix time"
// @Router       /api/file/download_task [get]
func downloadTaskHandler(ctx echo.Context) error {
	var msg download.Msg_Req_Download_Task
	res := &download.Msg_Resp_Download_Task{}

	if err := ctx.Bind(&msg); err != nil {
		res.MetaStatus(-1, "post data error")
		return ctx.JSON(http.StatusOK, res)
	}

	//check token
	err := CheckToken(ctx)
	if err != nil {
		res.MetaStatus(-2, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	//todo use a channel to accept task
	go download_mgr.StartDownloader(msg.Origin_url, msg.file_hash, download_mgr.SuccessCallback, download_mgr.FailedCallback)

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)

}
