package config

import (
	"github.com/meson-network/peer-node/configuration"
	"github.com/urfave/cli/v2"
)

//set your config params types
var stringConfParams = []string{}
var float64ConfParams = []string{}
var boolConfPrams = []string{}
var otherConf = []string{}

//get all config flags
func GetFlags() (allflags []cli.Flag) {
	allConfig := configuration.Config.AllSettings()
	for k, v := range allConfig {
		switch v.(type) {
		case string:
			stringConfParams = append(stringConfParams, k)
		case float64:
			float64ConfParams = append(float64ConfParams, k)
		case bool:
			boolConfPrams = append(boolConfPrams, k)
		}
	}

	for _, name := range stringConfParams {
		allflags = append(allflags, &cli.StringFlag{Name: name, Required: false})
	}

	for _, name := range float64ConfParams {
		allflags = append(allflags, &cli.Float64Flag{Name: name, Required: false})
	}

	for _, name := range boolConfPrams {
		allflags = append(allflags, &cli.BoolFlag{Name: name, Required: false})
	}

	//other custom flags
	allflags = append(allflags, &cli.StringFlag{Name: "addpath", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "removepath", Required: false})

	return
}
