package daemon_plugin

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/takama/daemon"
)

type Service struct {
	daemon.Daemon
}

var instanceMap = map[string]*Service{}

func GetInstance(service_name string) *Service {
	return instanceMap[service_name]
}

// Init a new instance.
//  If only need one instance, use empty name "". Use GetDefaultInstance() to get.
//  If you need several instance, run Init() with different <name>. Use GetInstance(<name>) to get.
func Init(name string) error {
	if name == "" {
		return errors.New("name can not be vacant ")
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("daemon instance <%s> has already been initialized", name)
	}

	kind := daemon.SystemDaemon
	if runtime.GOOS == "darwin" {
		kind = daemon.UserAgent
	}
	srv, err := daemon.New(name, name+":description", kind)
	if err != nil {
		return err
	}
	instanceMap[name] = &Service{srv}
	return nil
}
