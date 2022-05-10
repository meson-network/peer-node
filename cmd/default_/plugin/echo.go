package plugin

import (
	"errors"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/remote/cert"
	tool_errors "github.com/meson-network/peer-node/tools/errors"
)

func InitEchoServer() error {

	https_port, err := configuration.Config.GetInt("https_port", 443)
	if err != nil {
		return errors.New("https_port [int] in config error," + err.Error())
	}

	cert_mgr, cert_mgr_err := cert.GetCertMgr()
	if cert_mgr_err != nil {
		return cert_mgr_err
	}

	return echo_plugin.Init(echo_plugin.Config{Port: https_port, Tls: true, Crt_path: cert_mgr.Crt_path, Key_path: cert_mgr.Key_path},
		tool_errors.PanicHandler, basic.Logger)

}
