package version_mgr

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer_common/version"
)

const NodeVersion = "3.0.7"

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

func (v *VersionMgr) IsLatestVersion() (isLatestVersion bool, latestVersion string, downloadHost string, err error) {
	latestVersion, _, downloadHost, err = client.GetNodeVersionFromServer()
	if err != nil {
		return true, latestVersion, downloadHost, err
	}

	r := version.VersionCompare(v.CurrentVersion, latestVersion)
	if r >= 0 {
		return true, latestVersion, downloadHost, nil
	}
	return false, latestVersion, downloadHost, nil
}

func GetOSInfo() (arch string, osInfo string) {
	arch = "amd64"
	switch runtime.GOARCH {
	case "386":
		arch = "386"
	case "arm64":
		arch = "arm64"
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

func (v *VersionMgr) CheckUpdate() {
	isLatestVersion, latestVersion, downloadHost, _ := v.IsLatestVersion()
	if isLatestVersion {
		return
	}

	//download new version
	//need upgrade
	if v.AutoUpdateFiledTime > updateRetryTimeLimit || v.LastFailedTime > time.Now().UTC().Unix()-updateRetryIntervalSec {
		basic.Logger.Infoln("New version auto update failed, please update by manual.")
		return
	}
	basic.Logger.Infoln("New version detected, start to upgrade... ")

	//new version download url
	fileName := genFileName()
	downloadPath := "v" + latestVersion + "/" + fileName
	newVersionDownloadUrl := downloadHost + "/node/" + downloadPath
	basic.Logger.Debugln("new version download url", "url", newVersionDownloadUrl)

	//upgrade
	err := upgradeNewVersion(newVersionDownloadUrl)
	if err != nil {
		basic.Logger.Errorln("CheckUpdate DownloadNewVersion err:", err)
		v.LastFailedTime = time.Now().UTC().Unix()
		v.AutoUpdateFiledTime++
		return
	}

	//restart
	err = RestartNode()
	if err != nil {
		basic.Logger.Errorln("CheckUpdate RestartNode err:", err)
		v.LastFailedTime = time.Now().UTC().Unix()
		v.AutoUpdateFiledTime++
		return
	}
}

func upgradeNewVersion(downloadUrl string) error {
	//get
	response, err := http.Get(downloadUrl)
	if err != nil {
		basic.Logger.Errorln(" upgradeNewVersion get file url "+downloadUrl+" error", "err", err)
		return err
	}
	if response.Body == nil {
		basic.Logger.Errorln("upgradeNewVersion response body is null")
		return err
	}
	defer response.Body.Close()

	//unzip to temp folder
	tempFolder := "temp"
	tempFolder = path_util.ExE_Path(tempFolder)
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		basic.Logger.Errorln("upgradeNewVersion os.MkdirAll err", err, "filePath", tempFolder)
		return err
	}
	err = unzip(tempFolder, response.Body)
	if err != nil {
		basic.Logger.Errorln("upgradeNewVersion unzip err", err, "filePath", tempFolder)
		return err
	}

	//overwrite oldFile
	runningPath := path_util.ExE_PathStr()
	err = filepath.Walk(tempFolder, func(path string, info fs.FileInfo, err error) error {
		oldFile := filepath.Join(runningPath, info.Name())
		newFile := filepath.Join(tempFolder, info.Name())

		if info.IsDir() {
			err = os.MkdirAll(oldFile, 0777)
			if err != nil {
				basic.Logger.Errorln("upgradeNewVersion filepath.Walk os.MkdirAll err", err, "path:", oldFile)
				return err
			}
		}

		err = overwriteOldFile(newFile, oldFile)
		if err != nil {
			basic.Logger.Errorln("upgradeNewVersion filepath.Walk overwriteOldFile err", err, "path:", oldFile)
			return err
		}

		return nil
	})
	if err != nil {
		basic.Logger.Errorln("upgradeNewVersion overwrite old file filepath.Walk err", err)
		return err
	}
	os.RemoveAll(tempFolder)

	return nil
}
