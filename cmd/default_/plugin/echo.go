package plugin

import (
	"fmt"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/cert_mgr"
	tool_errors "github.com/meson-network/peer-node/tools/errors"
)

func InitEchoServer() error {

	toml_conf := conf.Get_config().Toml_config

	if toml_conf.Https_port == 0 || echo_plugin.IsForbiddenPort(toml_conf.Https_port) {
		return fmt.Errorf("[https_port] [%d] in config is forbidden, please use other port. 443 is recommended", toml_conf.Https_port)
	}

	return echo_plugin.Init(echo_plugin.Config{Port: toml_conf.Https_port, Tls: true, Crt_path: cert_mgr.GetInstance().Crt_path, Key_path: cert_mgr.GetInstance().Key_path},
		tool_errors.PanicHandler, basic.Logger)

}
