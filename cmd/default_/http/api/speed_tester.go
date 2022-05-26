package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/tools/http/api"
)

func config_speed_tester(httpServer *echo_plugin.EchoServer) {
	httpServer.GET("/api/speed_tester/pause/:second", pauseHandler, MidToken)
	httpServer.GET("/api/speed_tester/test", testHandler, MidToken)
}

// @Summary      /api/health
// @Description  health check
// @Tags         health
// @Produce      json
// @Success      200 {object} MSG_RESP_HEALTH "server unix time"
// @Router       /api/speed_tester/pause/{second} [get]
func pauseHandler(ctx echo.Context) error {
	res := &api.API_META_STATUS{}

	pauseTimeStr := ctx.Param("second")
	pauseTime, err := strconv.Atoi(pauseTimeStr)
	if err != nil {
		res.MetaStatus(-1, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}
	if pauseTime < 0 || pauseTime > 10 {
		pauseTime = 4
	}

	//check token
	err = CheckToken(ctx)
	if err != nil {
		res.MetaStatus(-2, err.Error())
		return ctx.JSON(http.StatusOK, res)
	}

	echo_plugin.GetInstance().SetPauseSeconds(int64(pauseTime))

	res.MetaStatus(1, "success")
	return ctx.JSON(http.StatusOK, res)
}

// @Summary      /api/health
// @Description  health check
// @Tags         health
// @Produce      json
// @Success      200 {object} MSG_RESP_HEALTH "server unix time"
// @Router       /api/speed_tester/pause/{second} [get]
func testHandler(ctx echo.Context) error {
	//check token
	err := CheckToken(ctx)
	if err != nil {
		return ctx.HTML(http.StatusUnauthorized, "")
	}

	//check file exist
	//todo send file

	return ctx.File("")
}
