package daemon_plugin

import (
	"fmt"
	"runtime"

	"github.com/takama/daemon"
)

const (
	// name of the service
	serviceName = "template"
	description = "app template"
)

type Service struct {
	daemon.Daemon
}

var instanceMap = map[string]*Service{}

func GetInstance() *Service {
	return instanceMap["default"]
}

func GetInstance_(name string) *Service {
	return instanceMap[name]
}

func Init() error {
	return Init_("default")
}

// Init a new instance.
//  If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//  If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init_(name string) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("daemon instance <%s> has already been initialized", name)
	}

	kind := daemon.SystemDaemon
	if runtime.GOOS == "darwin" {
		kind = daemon.UserAgent
	}
	srv, err := daemon.New(serviceName, description, kind)
	if err != nil {
		return err
	}
	instanceMap[name] = &Service{srv}
	return nil
}
