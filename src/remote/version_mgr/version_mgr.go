package version_mgr

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/version"
)

const NodeVersion = "3.0.0"

type VersionMgr struct {
	CurrentVersion string
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
		CurrentVersion: NodeVersion,
	}
	return nil
}

func (v *VersionMgr) GetVersion() string {
	return v.CurrentVersion
}

func (v *VersionMgr) getNodeVersionFromServer() (latestVersion string, allowVersion string, err error) {
	//check is there new version or not
	basic.Logger.Debugln("Check Version...")
	result := &version.Msg_Resp_NodeVersion{}
	err = api.Get_(client.EndPoint+"/api/node/version", client.Token, 30, result)
	if err != nil {
		return "", "", err
	}

	if result.Meta_status <= 0 {
		return "", "", errors.New(result.Meta_message)
	}

	return result.Latest_version, result.Allow_version, nil
}

func (v *VersionMgr) IsLatestVersion() (isLatestVersion bool, latestVersion string, err error) {
	latestVersion, _, err = v.getNodeVersionFromServer()
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
