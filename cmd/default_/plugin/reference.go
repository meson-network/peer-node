package plugin

import "github.com/meson-network/peer-node/plugin/reference_plugin"

//example 3 cache instance
func initReference() error {
	//default instance
	err := reference_plugin.Init()
	if err != nil {
		return err
	}

	return nil
}
