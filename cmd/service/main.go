package service

import (
	"os"

	"github.com/kardianos/service"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_"
	"github.com/meson-network/peer-node/src/precheck_config"
	"github.com/urfave/cli/v2"
)

type Program struct{}

func (p *Program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *Program) run() {
	// Do work here
	default_.StartDefault(nil)
}
func (p *Program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	//basic.Logger.Infoln("service will stop in 5 seconds...")
	//<-time.After(time.Second * 5)
	return nil
}

func RunServiceCmd(clictx *cli.Context, s service.Service) {

	exe_path, exe_path_err := os.Executable()
	if exe_path_err != nil {
		basic.Logger.Errorln(exe_path_err)
		return

	}

	basic.Logger.Debugln("exefile:" + exe_path + " to be service target")

	//check command
	subCmds := clictx.Command.Names()
	if len(subCmds) == 0 {
		basic.Logger.Fatalln("no sub command")
	}

	action := subCmds[0]

	//var status string
	switch action {
	case "install":
		//check config
		precheck_config.CheckConfig()
		err := s.Install()
		if err != nil {
			basic.Logger.Errorln("install service error:", err)
		} else {
			basic.Logger.Infoln("service installed")
		}

	case "remove":

		err := s.Uninstall()
		if err != nil {
			basic.Logger.Errorln("remove service error:", err)
		} else {
			basic.Logger.Infoln("service removed")
		}

	case "start":
		//check config
		precheck_config.CheckConfig()
		err := s.Start()
		if err != nil {
			basic.Logger.Errorln("start service error:", err)
		} else {
			basic.Logger.Infoln("service started")
		}

	case "stop":
		err := s.Stop()
		if err != nil {
			basic.Logger.Errorln("stop service error:", err)
		} else {
			basic.Logger.Infoln("service stopped")
		}

	case "restart":
		err := s.Restart()
		if err != nil {
			basic.Logger.Errorln("restart service error:", err)
		} else {
			basic.Logger.Infoln("service restarted")
		}

	case "status":
		status, err := s.Status()
		if err != nil {
			basic.Logger.Errorln(err)
		}
		switch status {
		case service.StatusRunning:
			basic.Logger.Infoln("service status:", "RUNNING")
		case service.StatusStopped:
			basic.Logger.Infoln("service status:", "STOPPED")
		default:
			basic.Logger.Infoln("service status:", "UNKNOWN")
		}
	default:
		basic.Logger.Debugln("no sub command")
		return
	}
}
