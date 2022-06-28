package precheck_config

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer_common/cdn_cache"
)

func PreCheckConfig() {
	toml_conf := conf.Get_config().Toml_config

	token := toml_conf.Token
	if len(token) != 24 {
		basic.Logger.Fatalln("token config error, please set correct token in config")
	}

	port := toml_conf.Https_port
	if port <= 0 || port > 65535 || echo_plugin.IsForbiddenPort(port) {
		basic.Logger.Fatalln("https_port cofig error, please set correct port in config")
	}

	cdnCacheSize := toml_conf.Cache.Size
	if cdnCacheSize < cdn_cache.MIN_CACHE_SIZE {
		basic.Logger.Fatalln("cache.size config error, minimum is 20, please set cache size in config")
	}

	if toml_conf.Storage.Enable == false {
		return
	}

	apiPort := toml_conf.Storage.Api_port
	if apiPort <= 0 || apiPort > 65535 {
		basic.Logger.Fatalln("storage.api_port error, please set correct port in config")
	}
	if apiPort == toml_conf.Https_port {
		basic.Logger.Fatalln("storage api port [%d] already used in https port", apiPort)
	}

	consolePort := toml_conf.Storage.Console_port
	if consolePort <= 0 || consolePort > 65535 {
		basic.Logger.Fatalln("storage.console_port error, please set correct port in config")
	}
	if consolePort == toml_conf.Https_port || consolePort == apiPort {
		basic.Logger.Fatalln("storage console port [%d] already used in https port or api port", consolePort)
	}

	password := toml_conf.Storage.Password
	if password == "" {
		basic.Logger.Fatalln("storage.password not exist in config")
	}
	if len(password) < 6 {
		basic.Logger.Fatalln("storage.password length can not less than 6")
	}

}
