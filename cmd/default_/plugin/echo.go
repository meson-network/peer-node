package plugin

import (
	"errors"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	tool_errors "github.com/meson-network/peer-node/tools/errors"
)

func initEchoServer() error {
	http_port, err := configuration.Config.GetInt("http_port", 8080)
	if err != nil {
		return errors.New("http_port [int] in config error," + err.Error())
	}

	return echo_plugin.Init(echo_plugin.Config{Port: http_port}, tool_errors.PanicHandler, basic.Logger)
}
