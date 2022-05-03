package service

import (
	"fmt"
	"os"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/daemon_plugin"
	"github.com/urfave/cli/v2"
)

func RunServiceCmd(clictx *cli.Context) {
	//check command
	subCmds := clictx.Command.Names()
	if len(subCmds) == 0 {
		basic.Logger.Fatalln("no sub command")
	}

	action := subCmds[0]
	err := daemon_plugin.Init()
	if err != nil {
		basic.Logger.Fatalln("init daemon service error:", err)
	}

	var status string
	var e error
	switch action {
	case "install":
		status, e = daemon_plugin.GetInstance().Install()
		basic.Logger.Debugln("cmd install")
	case "remove":
		daemon_plugin.GetInstance().Stop()
		status, e = daemon_plugin.GetInstance().Remove()
		basic.Logger.Debugln("cmd remove")
	case "start":
		status, e = daemon_plugin.GetInstance().Start()
		basic.Logger.Debugln("cmd start")
	case "stop":
		status, e = daemon_plugin.GetInstance().Stop()
		basic.Logger.Debugln("cmd stop")
	case "restart":
		daemon_plugin.GetInstance().Stop()
		status, e = daemon_plugin.GetInstance().Start()
		basic.Logger.Debugln("cmd restart")
	case "status":
		status, e = daemon_plugin.GetInstance().Status()
		basic.Logger.Debugln("cmd status")
	default:
		basic.Logger.Debugln("no sub command")
		return
	}

	if e != nil {
		fmt.Println(status, "\nError: ", e)
		os.Exit(1)
	}
	fmt.Println(status)
}
