package plugin

import (
	"errors"
	"fmt"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/cert_mgr"
	tool_errors "github.com/meson-network/peer-node/tools/errors"
)

func InitEchoServer() error {

	https_port, err := configuration.Config.GetInt("https_port", 443)
	if err != nil {
		return errors.New("https_port [int] in config error," + err.Error())
	}

	if echo_plugin.IsForbiddenPort(https_port) {
		return fmt.Errorf("https_port [%d] in config is forbidden, please use other port. 443 is recommended", https_port)
	}

	return echo_plugin.Init(echo_plugin.Config{Port: https_port, Tls: true, Crt_path: cert_mgr.GetInstance().Crt_path, Key_path: cert_mgr.GetInstance().Key_path},
		tool_errors.PanicHandler, basic.Logger)

}
