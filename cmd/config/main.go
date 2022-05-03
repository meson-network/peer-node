package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/meson-network/peer-node/configuration"
	"github.com/urfave/cli/v2"
)

func ConfigSetting(clictx *cli.Context) {

	configModify := false

	for _, v := range stringConfParams {
		if clictx.IsSet(v) {
			newValue := clictx.String(v)
			configuration.Config.Set(v, newValue)
			configModify = true
		}
	}

	for _, v := range float64ConfParams {
		if clictx.IsSet(v) {
			newValue := clictx.Float64(v)
			configuration.Config.Set(v, newValue)
			configModify = true
		}
	}

	for _, v := range boolConfPrams {
		if clictx.IsSet(v) {
			newValue := clictx.Bool(v)
			configuration.Config.Set(v, newValue)
			configModify = true
		}
	}

	if configModify {
		err := configuration.Config.WriteConfig()
		if err != nil {
			color.Red("config save error:", err)
			return
		}
		fmt.Println("config modified")
		//fmt.Println(configuration.Config.GetConfigAsString())
	}
}
