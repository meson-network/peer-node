package version_mgr

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer_common/version"
	"github.com/pelletier/go-toml"
)

const NodeVersion = "3.1.20"

const updateRetryIntervalSec = 12 * 3600
const updateRetryTimeLimit = 4

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

func (v *VersionMgr) CheckUpdate() {
	isLatestVersion, latestVersion, downloadHost, _ := v.IsLatestVersion()
	if isLatestVersion {
		return
	}

	// check main version
	// mainVersion := strings.Split(latestVersion, ".")[0]
	// currentMainVersion := strings.Split(NodeVersion, ".")[0]
	// if mainVersion != currentMainVersion {
	//	basic.Logger.Infoln("New version released, please download new version.")
	//	return
	// }

	if v.LastFailedTime > time.Now().UTC().Unix()-updateRetryIntervalSec {
		return
	}
	// one week later try again
	if v.AutoUpdateFiledTime >= updateRetryTimeLimit && v.LastFailedTime < time.Now().UTC().Unix()-5*3600*24 {
		v.AutoUpdateFiledTime = updateRetryTimeLimit - 1
	}

	if v.AutoUpdateFiledTime >= updateRetryTimeLimit {
		basic.Logger.Infoln("New version auto update failed, please update by manual.")
		return
	}

	// do upgrade
	basic.Logger.Infoln("New version detected, start to upgrade... ")

	// download new version
	// download url
	fileName := genFileName()
	downloadPath := "v" + latestVersion + "/" + fileName
	newVersionDownloadUrl := downloadHost + "/" + downloadPath
	basic.Logger.Debugln("new version download url", "url", newVersionDownloadUrl)

	// upgrade
	err := upgradeNewVersion(newVersionDownloadUrl)
	if err != nil {
		basic.Logger.Errorln("CheckUpdate DownloadNewVersion err:", err)
		v.LastFailedTime = time.Now().UTC().Unix()
		v.AutoUpdateFiledTime++
		return
	}

	// restart
	err = RestartNode()
	if err != nil {
		basic.Logger.Errorln("CheckUpdate RestartNode err:", err)
		v.LastFailedTime = time.Now().UTC().Unix()
		v.AutoUpdateFiledTime++
		return
	}
}

func upgradeNewVersion(downloadUrl string) error {
	// get
	// todo ignore tls???
	tt := &http.Transport{
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	requestClient := &http.Client{Transport: tt}
	response, err := requestClient.Get(downloadUrl)
	if err != nil {
		basic.Logger.Errorln(" upgradeNewVersion get file url "+downloadUrl+" error", "err", err)
		return err
	}
	if response.Body == nil {
		basic.Logger.Errorln("upgradeNewVersion response body is null")
		return err
	}
	defer response.Body.Close()

	// unzip to temp folder
	tempFolder := "upgradetemp"
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

	// merge config
	// handle default.toml
	// read new config
	newConfigFile := filepath.Join(tempFolder, "configs", "default.toml")
	newConfigTree, err := toml.LoadFile(newConfigFile)
	if err != nil {
		basic.Logger.Errorln("new config file toml.LoadFile err:", err, "path", newConfigFile)
		return err
	}

	// read upgrade_keep
	reserveKeyArray := []string{}
	reserveKey := newConfigTree.Get("upgrade_keep")
	if reserveKey != nil {
		for _, key := range reserveKey.([]interface{}) {
			reserveKeyArray = append(reserveKeyArray, key.(string))
		}
	}

	// read old config
	runningPath := path_util.ExE_PathStr()
	oldConfigFile := filepath.Join(runningPath, "configs", "default.toml")
	oldConfigTree, err := toml.LoadFile(oldConfigFile)
	if err != nil {
		basic.Logger.Errorln("old config file toml.LoadFile err:", err, "filePath", oldConfigFile)
		return err
	}

	// merge config content
	config := mergeConfig(oldConfigTree, newConfigTree, reserveKeyArray)
	content, err := config.Marshal()
	if err != nil {
		basic.Logger.Errorln("mergeConfig Marshal err:", err)
		return err
	}
	err = os.WriteFile(newConfigFile, content, 0777)
	if err != nil {
		basic.Logger.Errorln("mergeConfig write newConfig file err:", err)
		return err
	}

	// overwrite oldFile
	err = filepath.WalkDir(tempFolder, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "upgradetemp" {
			return nil
		}

		relPath, err := filepath.Rel(tempFolder, path)
		if err != nil {
			basic.Logger.Errorln("upgradeNewVersion filepath.Walk filepath.Rel err", err, "path:", path)
			return err
		}
		oldFile := filepath.Join(runningPath, relPath)
		newFile := path

		if d.IsDir() {
			err = os.MkdirAll(oldFile, 0777)
			if err != nil {
				basic.Logger.Errorln("upgradeNewVersion filepath.Walk os.MkdirAll err", err, "path:", oldFile)
				return err
			}
			return nil
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
