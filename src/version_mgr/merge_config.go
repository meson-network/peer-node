package version_mgr

import (
	"github.com/pelletier/go-toml"
)

func mergeConfig(oldConfigTree *toml.Tree, newConfig *toml.Tree, reserveKeys []string) (config *toml.Tree) {
	for _, key := range reserveKeys {
		value := oldConfigTree.Get(key)
		if value == nil {
			continue
		}
		newConfig.Set(key, value)
	}
	return newConfig
}
