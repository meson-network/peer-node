package daemon_plugin

import (
	"errors"
	"fmt"

	"github.com/kardianos/service"
	"github.com/meson-network/peer-node/basic"
)

var instanceMap = map[string]service.Service{}

func GetInstance(service_name string) service.Service {
	return instanceMap[service_name]
}

// Init a new instance.
//  If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//  If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init(name string, program service.Interface) error {
	if name == "" {
		return errors.New("name can not be vacant ")
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("daemon instance <%s> has already been initialized", name)
	}

	svcConfig := &service.Config{
		Name:        name,
		DisplayName: name,
		Description: name + ":description",

		Option: map[string]interface{}{
			"OnFailure":              "restart",
			"OnFailureDelayDuration": "15s",
			"SystemdScript":          systemdScript,
			"Restart":                "on-failure",
		},
	}

	s, err := service.New(program, svcConfig)
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	instanceMap[name] = s
	return nil
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
StandardOutput=file:/var/log/{{.Name}}.out
StandardError=file:/var/log/{{.Name}}.err
{{- end}}
{{if gt .LimitNOFILE -1 }}LimitNOFILE={{.LimitNOFILE}}{{end}}
{{if .Restart}}Restart={{.Restart}}{{end}}
{{if .SuccessExitStatus}}SuccessExitStatus={{.SuccessExitStatus}}{{end}}
RestartSec=60
EnvironmentFile=-/etc/sysconfig/{{.Name}}

[Install]
WantedBy=multi-user.target
`
