package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http/api"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
)

//httpServer example
func StartDefaultHttpSever() {
	httpServer := echo_plugin.GetInstance()
	api.ConfigApi(httpServer)
	api.DeclareApi(httpServer)

	//for handling storage
	httpServer.GET("/*", func(ctx echo.Context) error {
		//storage_mgr.GetInstance()
		return ctx.HTML(http.StatusOK, "default")
	})

	err := httpServer.Start()
	if err != nil {
		basic.Logger.Fatalln(err)
	}
}

func CheckDefaultHttpServerStarted() bool {
	return echo_plugin.GetInstance().CheckStarted()
}

func ServerReloadCert() error {
	return echo_plugin.GetInstance().ReloadCert()
}
