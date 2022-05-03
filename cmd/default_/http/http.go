package http

import (
	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http/api"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
)

//httpServer example
func StartDefaultHttpSever() {
	httpServer := echo_plugin.GetInstance()
	api.ConfigApi(httpServer)
	api.DeclareApi(httpServer)

	//static
	conf_http_static_dir, sd_err := configuration.Config.GetString("http_static_dir", "")
	if sd_err == nil && conf_http_static_dir != "" {
		h_s_d, err := path_util.SmartExistPath(conf_http_static_dir)
		if err == nil {
			httpServer.Static("/", h_s_d)
			basic.Logger.Infoln("http static folder:", h_s_d)
		}
	}

	err := httpServer.Start()
	if err != nil {
		basic.Logger.Fatalln(err)
	}
}

func CheckDefaultHttpServerStarted() bool {
	return echo_plugin.GetInstance().CheckStarted()
}
