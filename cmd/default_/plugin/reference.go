package plugin

import "github.com/meson-network/peer-node/plugin/reference_plugin"

//example 3 cache instance
func initReference() error {
	//default instance
	err := reference_plugin.Init()
	if err != nil {
		return err
	}

	// cache1 instance
	err = reference_plugin.Init_("ref1")
	if err != nil {
		return err
	}

	// cache2 instance
	err = reference_plugin.Init_("ref2")
	if err != nil {
		return err
	}

	return nil
}
