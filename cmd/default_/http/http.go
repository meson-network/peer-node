package http

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http/api"
	"github.com/meson-network/peer-node/cmd/default_/http/file_request"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
)

//httpServer example
func StartDefaultHttpSever() {
	httpServer := echo_plugin.GetInstance()
	api.ConfigApi(httpServer)
	api.DeclareApi(httpServer)

	//for handling private storage
	//httpServer.GET("/_personal_/*", func(ctx echo.Context) error {
	//	//storage_mgr.GetInstance()
	//	return ctx.HTML(http.StatusOK, "personal data")
	//})

	//for handling cache file request
	// https://spec00-xxsdfsdffsdf-06-pzname.xxx.com/path1/path2/path3/1.jpg
	file_request.HandleFileRequest(httpServer)

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
