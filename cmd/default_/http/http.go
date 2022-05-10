package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http/api"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/file_mgr"
)

//httpServer example
func StartDefaultHttpSever() {
	httpServer := echo_plugin.GetInstance()
	api.ConfigApi(httpServer)
	api.DeclareApi(httpServer)

	//for handling private storage
	httpServer.GET("/_personal_/*", func(ctx echo.Context) error {
		//storage_mgr.GetInstance()
		return ctx.HTML(http.StatusOK, "personal data")
	})

	//for handling public storage
	httpServer.GET("/*", func(ctx echo.Context) error {
		// access_token := ctx.Request().Header.Get("access_token")
		// if access_token == "" {
		// 	return ctx.HTML(http.StatusOK, "request is forbidden")
		// }

		basic.Logger.Infoln(ctx.Request().URL)
		basic.Logger.Infoln(file_mgr.UrlToPublicFileHash(ctx.Request().RequestURI))
		basic.Logger.Infoln(file_mgr.UrlToPublicFileRelPath(ctx.Request().RequestURI))

		file_abs, file_abs_err := file_mgr.RequestPublicFile(file_mgr.UrlToPublicFileHash(ctx.Request().RequestURI))

		basic.Logger.Infoln(file_abs)
		if file_abs_err != nil {
			basic.Logger.Debugln(file_abs_err)
			return ctx.HTML(404, "file not found")
		}

		return ctx.File(file_abs)
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
