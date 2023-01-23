package precheck_config

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer_common/cdn_cache"
)

func PreCheckConfig() {
	toml_conf := conf.Get_config().Toml_config

	port := toml_conf.Https_port
	if port <= 0 || port > 65535 || echo_plugin.IsForbiddenPort(port) {
		basic.Logger.Fatalln("https_port cofig error, please set correct port in config")
	}

	cdnCacheSize := toml_conf.Cache.Size
	if cdnCacheSize < cdn_cache.MIN_CACHE_SIZE {
		basic.Logger.Fatalln("cache.size config error, minimum is 20, please set cache size in config")
	}

}
