//go:build linux
// +build linux

package version_mgr

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
)

func (v *VersionMgr) CheckUpdate() {
	isLatestVersion, latestVersion, _ := v.IsLatestVersion()
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
	//check arch and os
	arch, osInfo := GetOSInfo()
	// 'https://dashboard.meson.network/static_assets/node/v0.1.2/meson-darwin-amd64.tar.gz'
	fileName := "meson" + "-" + osInfo + "-" + arch + ".tar.gz"
	downloadPath := "v" + latestVersion + "/" + fileName
	newVersionDownloadUrl := "https://dashboard.meson.network/static_assets/node/" + downloadPath
	basic.Logger.Debugln("new version download url", "url", newVersionDownloadUrl)

	err := DownloadNewVersion(newVersionDownloadUrl)
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

func DownloadNewVersion(downloadUrl string) error {
	//get
	response, err := http.Get(downloadUrl)
	if err != nil {
		basic.Logger.Errorln(" DownloadNewVersion get file url "+downloadUrl+" error", "err", err)
		return err
	}

	//defer file.Close()
	if response.Body == nil {
		return errors.New("DownloadNewVersion body is null")
	}
	defer response.Body.Close()

	// gzip read
	gr, err := gzip.NewReader(response.Body)
	if err != nil {
		basic.Logger.Errorln("DownloadNewVersion gzip read new version file error", err)
		return err
	}
	defer gr.Close()
	// tar read
	tr := tar.NewReader(gr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			basic.Logger.Errorln("DownloadNewVersion unzip new version file error", err)
			return err
		}

		arr := strings.Split(h.Name, "/")
		nameArr := []string{}
		for _, v := range arr {
			if v != "" {
				nameArr = append(nameArr, v)
			}
		}
		if len(nameArr) <= 1 {
			continue
		}
		name := filepath.Join(nameArr[1:]...)

		//skip config folder and pro.json
		if name == "configs" || name == "configs/pro.json" {
			continue
		}

		filePath := path_util.ExE_Path(name)
		if h.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, 0777)
			if err != nil {
				basic.Logger.Errorln("DownloadNewVersion os.MkdirAll err", err, "filePath", filePath)
				return err
			}
			continue
		}

		content, err := ioutil.ReadAll(tr)
		if err != nil {
			basic.Logger.Errorln("DownloadNewVersion ioutil.ReadAll err", err, "filePath", filePath)
			return err
		}

		err = os.Remove(filePath)
		if err != nil {
			basic.Logger.Errorln("Error remove old file", filePath, "err:", err)
			fmt.Println(err)
			return err
		}
		err = ioutil.WriteFile(filePath, content, 777)
		if err != nil {
			basic.Logger.Errorln("Error creating", filePath, "err:", err)
			return err
		}
		os.Chmod(filePath, 0777)

		//err = os.WriteFile(filePath, content, 0777)
		//if err != nil {
		//	basic.Logger.Errorln("DownloadNewVersion os.WriteFile err:", err, "filePath", filePath)
		//	return err
		//}
	}

	return nil
}

func RestartNode() error {
	basic.Logger.Debugln("peer node restart cmd")

	absPath, err := path_util.SmartExistPath("./meson")
	if err != nil {
		basic.Logger.Errorln("RestartNode path_util.SmartExistPath err:", err)
		return err
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("sudo %s service restart", absPath))
	err = cmd.Run()
	if err != nil {
		basic.Logger.Errorln("restart peer node error:", err)
		return err
	}
	return nil
}
