package config

import (
	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/urfave/cli/v2"
)

func Cli_get_flags() []cli.Flag {

	allflags := []cli.Flag{}
	allflags = append(allflags, &cli.StringFlag{Name: "log_level", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "token", Required: false})
	allflags = append(allflags, &cli.IntFlag{Name: "https_port", Required: false})

	allflags = append(allflags, &cli.IntFlag{Name: "cache.size", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "cache.folder", Required: false})

	allflags = append(allflags, &cli.BoolFlag{Name: "storage.enable", Required: false})
	allflags = append(allflags, &cli.IntFlag{Name: "storage.api_port", Required: false})
	allflags = append(allflags, &cli.IntFlag{Name: "storage.console_port", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "storage.folder", Required: false})
	allflags = append(allflags, &cli.StringFlag{Name: "storage.password", Required: false})

	return allflags
}

func Cli_set_config(clictx *cli.Context) {
	config := conf.Get_config()

	if clictx.IsSet("log_level") {
		config.Toml_config.Log_level = clictx.String("log_level")
	}
	if clictx.IsSet("token") {
		config.Toml_config.Token = clictx.String("token")
	}
	if clictx.IsSet("https_port") {
		config.Toml_config.Https_port = clictx.Int("https_port")
	}

	//cache config
	if clictx.IsSet("cache.size") {
		config.Toml_config.Cache.Size = clictx.Int("cache.size")
	}
	if clictx.IsSet("cache.folder") {
		config.Toml_config.Cache.Folder = clictx.String("cache.folder")
	}

	//storage
	if clictx.IsSet("storage.enable") {
		config.Toml_config.Storage.Enable = clictx.Bool("token")
	}
	if clictx.IsSet("storage.api_port") {
		config.Toml_config.Storage.Api_port = clictx.Int("storage.api_port")
	}
	if clictx.IsSet("storage.console_port") {
		config.Toml_config.Storage.Console_port = clictx.Int("storage.console_port")
	}
	if clictx.IsSet("storage.folder") {
		config.Toml_config.Storage.Folder = clictx.String("storage.folder")
	}
	if clictx.IsSet("storage.password") {
		config.Toml_config.Storage.Password = clictx.String("storage.password")
	}

	err := config.Save_config()
	if err != nil {
		color.Red("save config error:", err)
	}
}
