package config

import (
	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer_common/cdn_cache"
	"github.com/urfave/cli/v2"
)

func Cli_get_flags() []cli.Flag {

	allflags := []cli.Flag{}
	allflags = append(allflags, &cli.StringFlag{Name: "log.level", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "token", Required: false})
	allflags = append(allflags, &cli.IntFlag{Name: "https_port", Required: false})

	allflags = append(allflags, &cli.IntFlag{Name: "cache.size", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "cache.folder", Required: false})

	return allflags
}

func Cli_set_config(clictx *cli.Context) {
	config := conf.Get_config()

	if clictx.IsSet("log.level") {
		config.Toml_config.Log.Level = clictx.String("log.level")
	}
	if clictx.IsSet("token") {
		token := clictx.String("token")
		if len(token) != 24 {
			basic.Logger.Fatalln("token format error,token length should be 24")
		}
		config.Toml_config.Token = token
	}

	if clictx.IsSet("https_port") {
		port := clictx.Int("https_port")
		if port <= 0 || port > 65535 || echo_plugin.IsForbiddenPort(port) {
			basic.Logger.Fatalln("https_port config error,", port, "is not allowed")
		}

		config.Toml_config.Https_port = port
	}

	//cache config
	if clictx.IsSet("cache.size") {
		size := clictx.Int("cache.size")
		if size < cdn_cache.MIN_CACHE_SIZE {
			basic.Logger.Fatalln("cache.size config error,minimum is 20")
		}
		config.Toml_config.Cache.Size = size
	}
	if clictx.IsSet("cache.folder") {
		config.Toml_config.Cache.Folder = clictx.String("cache.folder")
	}

	err := config.Save_config()
	if err != nil {
		color.Red("save config error:", err)
	}
}
