package plugin

import (
	"errors"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	tool_errors "github.com/meson-network/peer-node/tools/errors"
)

func InitEchoServer() error {

	https_port, err := configuration.Config.GetInt("https_port", 443)
	if err != nil {
		return errors.New("https_port [int] in config error," + err.Error())
	}

	crt, err := configuration.Config.GetString("https_crt_path", "")
	if err != nil || crt == "" {
		return errors.New("https_crt_path [string] in config.json err")
	}

	key, err := configuration.Config.GetString("https_key_path", "")
	if err != nil || key == "" {
		return errors.New("https_key_path [string] in config.json err")
	}

	crt_path, cert_path_err := path_util.SmartExistPath(crt)
	if cert_path_err != nil {
		return errors.New("https crt file path error," + cert_path_err.Error())
	}

	key_path, key_path_err := path_util.SmartExistPath(key)
	if cert_path_err != nil {
		return errors.New("https key file path error," + key_path_err.Error())
	}

	return echo_plugin.Init(echo_plugin.Config{Port: https_port, Tls: true, Crt_path: crt_path, Key_path: key_path},
		tool_errors.PanicHandler, basic.Logger)

}
