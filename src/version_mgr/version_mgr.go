package version_mgr

import (
	"fmt"
	"runtime"

	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer_common/version"
)

const NodeVersion = "3.0.3"

const updateRetryIntervalSec = 12 * 3600
const updateRetryTimeLimit = 7

type VersionMgr struct {
	CurrentVersion      string
	AutoUpdateFiledTime int
	LastFailedTime      int64
}

var instanceMap = map[string]*VersionMgr{}

func GetInstance() *VersionMgr {
	return instanceMap["default"]
}

func GetInstance_(name string) *VersionMgr {
	return instanceMap[name]
}

func Init() error {
	return Init_("default")
}

func Init_(name string) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("echo server instance <%s> has already initialized", name)
	}
	instanceMap[name] = &VersionMgr{
		CurrentVersion:      NodeVersion,
		AutoUpdateFiledTime: 0,
		LastFailedTime:      0,
	}
	return nil
}

func (v *VersionMgr) GetVersion() string {
	return v.CurrentVersion
}

func (v *VersionMgr) IsLatestVersion() (isLatestVersion bool, latestVersion string, err error) {
	latestVersion, _, err = client.GetNodeVersionFromServer()
	if err != nil {
		return true, latestVersion, err
	}

	r := version.VersionCompare(v.CurrentVersion, latestVersion)
	if r >= 0 {
		return true, latestVersion, nil
	}
	return false, latestVersion, nil
}

func GetOSInfo() (arch string, osInfo string) {
	arch = "amd64"
	switch runtime.GOARCH {
	case "386":
		arch = "386"
	}

	osInfo = "linux"
	switch runtime.GOOS {
	case "windows":
		osInfo = "windows"
	case "darwin":
		osInfo = "darwin"
	}

	return arch, osInfo
}
