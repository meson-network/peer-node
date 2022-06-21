package cmd

import (
	"github.com/kardianos/service"
	"github.com/meson-network/peer-node/basic"
)

/////////service/////////////

func OS_service_start(name string, action string, exe_func func()) {
	os_service_conf := &service.Config{
		Name:        name,
		DisplayName: name,
		Description: name + ":description",

		Option: map[string]interface{}{
			"OnFailure":              "restart",
			"OnFailureDelayDuration": "15s",
			"SystemdScript":          systemdScript,
			"Restart":                "on-failure", // or use "always"
		},
	}

	oss, err := service.New(OS_service_program{Exe_func: exe_func}, os_service_conf)
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	os_service_run(&oss, action)
}

type OS_service_program struct {
	Exe_func func()
}

func (p OS_service_program) Start(s service.Service) error {
	if p.Exe_func != nil {
		go p.Exe_func()
	}
	return nil
}

func (p OS_service_program) Stop(s service.Service) error {
	return nil
}

func os_service_run(s *service.Service, action string) {
	switch action {
	case "install":
		err := (*s).Install()
		if err != nil {
			basic.Logger.Fatalln("install service error:", err)
		} else {
			basic.Logger.Infoln("service installed")
		}
	case "remove":
		err := (*s).Uninstall()
		if err != nil {
			basic.Logger.Fatalln("remove service error:", err)
		} else {
			basic.Logger.Infoln("service removed")
		}
	case "start":
		err := (*s).Start()
		if err != nil {
			basic.Logger.Fatalln("start service error:", err)
		} else {
			basic.Logger.Infoln("service started")
		}
	case "run":
		err := (*s).Run()
		if err != nil {
			basic.Logger.Fatalln("run service error:", err)
		} else {
			basic.Logger.Infoln("service run")
		}
	case "stop":
		err := (*s).Stop()
		if err != nil {
			basic.Logger.Fatalln("stop service error:", err)
		} else {
			basic.Logger.Infoln("service stopped")
		}
	case "restart":
		err := (*s).Restart()
		if err != nil {
			basic.Logger.Fatalln("restart service error:", err)
		} else {
			basic.Logger.Infoln("service restarted")
		}
	case "status":
		status, err := (*s).Status()
		if err != nil {
			basic.Logger.Fatalln(err)
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
		basic.Logger.Warnln("no sub command")
		return
	}
}

const systemdScript = `[Unit]
Description={{.Description}}
ConditionFileIsExecutable={{.Path|cmdEscape}}
{{range $i, $dep := .Dependencies}} 
{{$dep}} {{end}}

[Service]
StartLimitInterval=15
StartLimitBurst=15
ExecStart={{.Path|cmdEscape}}{{range .Arguments}} {{.|cmd}}{{end}}
{{if .ChRoot}}RootDirectory={{.ChRoot|cmd}}{{end}}
{{if .WorkingDirectory}}WorkingDirectory={{.WorkingDirectory|cmdEscape}}{{end}}
{{if .UserName}}User={{.UserName}}{{end}}
{{if .ReloadSignal}}ExecReload=/bin/kill -{{.ReloadSignal}} "$MAINPID"{{end}}
{{if .PIDFile}}PIDFile={{.PIDFile|cmd}}{{end}}
{{if and .LogOutput .HasOutputFileSupport -}}
StandardOutput=null
StandardError=null
{{- end}}
{{if gt .LimitNOFILE -1 }}LimitNOFILE={{.LimitNOFILE}}{{end}}
{{if .Restart}}Restart={{.Restart}}{{end}}
{{if .SuccessExitStatus}}SuccessExitStatus={{.SuccessExitStatus}}{{end}}
RestartSec=60
EnvironmentFile=-/etc/sysconfig/{{.Name}}

[Install]
WantedBy=multi-user.target
`
