package plugin

import "github.com/meson-network/peer-node/plugin/hub_plugin"

func initHub() error {
	return hub_plugin.Init()
}
