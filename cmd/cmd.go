package cmd

import (
	"fmt"
	"os"
	"strings"

	ilog "github.com/coreservice-io/log"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/cmd/config"
	"github.com/meson-network/peer-node/cmd/default_"
	"github.com/meson-network/peer-node/cmd/log"
	"github.com/meson-network/peer-node/src/precheck_config"
	"github.com/urfave/cli/v2"
)

const daemon_name = "meson-node"

const CMD_NAME_DEFAULT = "default"
const CMD_NAME_GEN_API = "gen_api"
const CMD_NAME_LOG = "log"
const CMD_NAME_SERVICE = "service"
const CMD_NAME_CONFIG = "config"

////////config to do cmd ///////////
func ConfigCmd() *cli.App {

	//////////init config/////////////
	toml_conf_path := "configs/default.toml"

	real_args := []string{}
	for _, arg := range os.Args {
		arg_lower := strings.ToLower(arg)
		if strings.HasPrefix(arg_lower, "-conf=") || strings.HasPrefix(arg_lower, "--conf=") {

			toml_target := strings.Trim(arg_lower, "-conf=")
			toml_target = strings.Trim(toml_target, "--conf=")
			toml_conf_path = "configs/" + toml_target + ".toml"
			fmt.Println("toml_conf_path", toml_conf_path)
			continue
		}
		real_args = append(real_args, arg)
	}

	os.Args = real_args

	conf_err := conf.Init_config(toml_conf_path)
	if conf_err != nil {
		basic.Logger.Fatalln("config err", conf_err)
	}

	configuration := conf.Get_config()

	/////set loglevel//////
	basic.Logger.SetLevel(ilog.ParseLogLevel(configuration.Toml_config.Log_level))
	////////////////////////////////

	return &cli.App{
		Action: func(clictx *cli.Context) error {
			OS_service_start(daemon_name, "run", func() {
				default_.StartDefault(clictx)
			})
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
			//{
			//	Name:  CMD_NAME_GEN_API,
			//	Usage: "api command",
			//	Action: func(clictx *cli.Context) error {
			//		api.Gen_Api_Docs()
			//		return nil
			//	},
			//},
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
							configs, _ := conf.Get_config().Read_config_file()
							fmt.Println(configs)
							fmt.Println("======== end  of  config ========")
							return nil
						},
					},
					//set config
					{
						Name:  "set",
						Usage: "set config",
						Flags: append(config.Cli_get_flags(), &cli.StringFlag{Name: "config", Required: false}),
						Action: func(clictx *cli.Context) error {
							config.Cli_set_config(clictx)
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
						Usage: "install service",
						Action: func(clictx *cli.Context) error {
							OS_service_start(daemon_name, "install", nil)
							return nil
						},
					},
					//service remove
					{
						Name:  "remove",
						Usage: "remove service",
						Action: func(clictx *cli.Context) error {
							OS_service_start(daemon_name, "remove", nil)
							return nil
						},
					},
					//service start
					{
						Name:  "start",
						Usage: "run",
						Action: func(clictx *cli.Context) error {
							//check config
							precheck_config.PreCheckConfig()
							OS_service_start(daemon_name, "start", nil)
							return nil
						},
					},
					//service stop
					{
						Name:  "stop",
						Usage: "stop",
						Action: func(clictx *cli.Context) error {
							OS_service_start(daemon_name, "stop", nil)
							return nil
						},
					},
					//service restart
					{
						Name:  "restart",
						Usage: "restart",
						Action: func(clictx *cli.Context) error {
							//check config
							precheck_config.PreCheckConfig()
							OS_service_start(daemon_name, "restart", nil)
							return nil
						},
					},
					//service status
					{
						Name:  "status",
						Usage: "show process status",
						Action: func(clictx *cli.Context) error {
							OS_service_start(daemon_name, "status", nil)
							return nil
						},
					},
				},
			},
		},
	}
}
