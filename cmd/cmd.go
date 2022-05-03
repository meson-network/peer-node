package cmd

import (
	"fmt"
	"os"
	"strings"

	ilog "github.com/coreservice-io/log"
	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/config"
	"github.com/meson-network/peer-node/cmd/default_"
	"github.com/meson-network/peer-node/cmd/default_/http/api"
	"github.com/meson-network/peer-node/cmd/log"
	"github.com/meson-network/peer-node/cmd/service"
	"github.com/meson-network/peer-node/configuration"
	"github.com/urfave/cli/v2"
)

const CMD_NAME_DEFAULT = "default"
const CMD_NAME_GEN_API = "gen_api"
const CMD_NAME_LOG = "log"
const CMD_NAME_SERVICE = "service"
const CMD_NAME_CONFIG = "config"

////////config to do cmd ///////////
func ConfigCmd() *cli.App {
	//check is dev or pro
	isDev := false
	confShow := false
	real_args := []string{}

	for _, arg := range os.Args {

		s := strings.ToLower(arg)
		if strings.Contains(s, "-mode=dev") || strings.Contains(s, "--mode=dev") {
			isDev = true
			continue
		}

		if strings.Contains(s, "-mode=pro") || strings.Contains(s, "--mode=pro") {
			isDev = false
			continue
		}

		if strings.Contains(s, "-conf=show") || strings.Contains(s, "--conf=show") {
			confShow = true
			continue
		}

		if strings.Contains(s, "-conf=hide") || strings.Contains(s, "--conf=hide") {
			confShow = false
			continue
		}

		real_args = append(real_args, arg)
	}

	os.Args = real_args

	conferr := iniConfig(isDev, confShow)
	if conferr != nil {
		basic.Logger.Panicln(conferr)
	}

	return &cli.App{
		Action: func(clictx *cli.Context) error {
			default_.StartDefault(clictx)
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  CMD_NAME_LOG,
				Usage: "print all logs",
				Flags: log.GetFlags(),
				Action: func(clictx *cli.Context) error {
					log.StartLog(clictx)
					return nil
				},
			},
			{
				Name:  CMD_NAME_GEN_API,
				Usage: "api command",
				Action: func(clictx *cli.Context) error {
					api.Gen_Api_Docs()
					return nil
				},
			},
			{
				Name:  CMD_NAME_CONFIG,
				Usage: "config command",
				Subcommands: []*cli.Command{
					//show config
					{
						Name:  "show",
						Usage: "show configs",
						Action: func(clictx *cli.Context) error {
							fmt.Println("======== start of config ========")
							configs, _ := configuration.Config.GetConfigAsString()
							fmt.Println(configs)
							fmt.Println("======== end  of  config ========")
							return nil
						},
					},
					//set config
					{
						Name:  "set",
						Usage: "set config",
						Flags: config.GetFlags(),
						Action: func(clictx *cli.Context) error {
							config.ConfigSetting(clictx)
							return nil
						},
					},
				},
			},
			{
				Name:  CMD_NAME_SERVICE,
				Usage: "service command",
				Subcommands: []*cli.Command{
					//service install
					{
						Name:  "install",
						Usage: "install meson node in service",
						Action: func(clictx *cli.Context) error {
							service.RunServiceCmd(clictx)
							return nil
						},
					},
					//service remove
					{
						Name:  "remove",
						Usage: "remove meson node from service",
						Action: func(clictx *cli.Context) error {
							service.RunServiceCmd(clictx)
							return nil
						},
					},
					//service start
					{
						Name:  "start",
						Usage: "run",
						Action: func(clictx *cli.Context) error {
							service.RunServiceCmd(clictx)
							return nil
						},
					},
					//service stop
					{
						Name:  "stop",
						Usage: "stop",
						Action: func(clictx *cli.Context) error {
							service.RunServiceCmd(clictx)
							return nil
						},
					},
					//service restart
					{
						Name:  "restart",
						Usage: "restart",
						Action: func(clictx *cli.Context) error {
							service.RunServiceCmd(clictx)
							return nil
						},
					},
					//service status
					{
						Name:  "status",
						Usage: "show process status",
						Action: func(clictx *cli.Context) error {
							service.RunServiceCmd(clictx)
							return nil
						},
					},
				},
			},
		},
	}
}

////////end config to do app ///////////
func readDefaultConfig(isDev bool, confShow bool) (*configuration.VConfig, string, error) {
	var defaultConfigPath string
	var err error
	if isDev {
		basic.Logger.Infoln("======== using dev mode ========")
		defaultConfigPath, err = path_util.SmartExistPath("configs/dev.json")
		if err != nil {
			basic.Logger.Errorln("no dev.json under /configs folder , use --mode=pro to run pro mode")
			return nil, "", err
		}
	} else {
		basic.Logger.Infoln("======== using pro mode ========")
		defaultConfigPath, err = path_util.SmartExistPath("configs/pro.json")
		if err != nil {
			basic.Logger.Errorln("no pro.json under /configs folder , use --mode=dev to run dev mode")
			return nil, "", err
		}
	}

	if confShow {
		basic.Logger.Infoln("using config:", defaultConfigPath)
	}

	config, err := configuration.ReadConfig(defaultConfigPath)
	if err != nil {
		basic.Logger.Errorln("config err", err)
		return nil, "", err
	}

	return config, defaultConfigPath, nil
}

func iniConfig(isDev bool, confShow bool) error {
	//path_util.ExEPathPrintln()
	////read default config
	config, _, err := readDefaultConfig(isDev, confShow)
	if err != nil {
		return err
	}

	configuration.Config = config
	logerr := setLoggerLevel()
	if logerr != nil {
		return logerr
	}

	if confShow {
		basic.Logger.Infoln("======== start of config ========")
		configs, _ := config.GetConfigAsString()
		basic.Logger.Infoln(configs)
		basic.Logger.Infoln("======== end  of  config ========")
	}

	return nil
}

func setLoggerLevel() error {
	logLevel := "INFO"
	if configuration.Config != nil {
		var err error
		logLevel, err = configuration.Config.GetString("local_log_level", "INFO")
		if err != nil {
			return err
		}
	}

	l := ilog.ParseLogLevel(logLevel)
	basic.Logger.SetLevel(l)
	return nil
}
