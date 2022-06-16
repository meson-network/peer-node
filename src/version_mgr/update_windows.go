//go:build windows
// +build windows

package version_mgr

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
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
	isLatestVersion, latestVersion, downloadHost, _ := v.IsLatestVersion()
	if isLatestVersion {
		return
	}

	basic.Logger.Infoln("New version detected, please update by manual.")
	return

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

	// 'https://meson.network/static/terminal/v0.1.2/meson-windows-amd64.zip'
	fileName := "meson" + "-" + osInfo + "-" + arch + ".zip"
	downloadPath := "v" + latestVersion + "/" + fileName
	newVersionDownloadUrl := downloadHost + "/static_assets/node/" + downloadPath
	basic.Logger.Debugln("new version download url", "url", newVersionDownloadUrl)

	err := DownloadNewVersion(newVersionDownloadUrl)
	if err != nil {
		basic.Logger.Errorln("CheckUpdate DownloadNewVersion err:", err)
		v.LastFailedTime = time.Now().UTC().Unix()
		v.AutoUpdateFiledTime++
		return
	}

	//restart
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

	if response == nil {
		return errors.New("response is nil")
	}

	if response.Body == nil {
		return errors.New("DownloadNewVersion body is null")
	}
	defer response.Body.Close()

	//get content
	tempContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		basic.Logger.Errorln("DownloadNewVersion ioutil.ReadAll err", err)
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(tempContent), int64(len(tempContent)))
	if err != nil {
		basic.Logger.Errorln("DownloadNewVersion zip.OpenReader err", err)
		return err
	}

	for _, f := range zipReader.File {
		arr := strings.Split(f.Name, "/")
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
		if name == "configs" || name == "configs/default.toml" {
			continue
		}

		filePath := path_util.ExE_Path(name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(filePath, 0777)
			if err != nil {
				basic.Logger.Errorln("DownloadNewVersion os.MkdirAll err", err, "filePath", filePath)
				return err
			}
			continue
		}

		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		content, err := ioutil.ReadAll(inFile)
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

	}

	return nil
}

func RestartNode() error {
	basic.Logger.Debugln("peer node restart cmd")

	absPath, exist, err := path_util.SmartPathExist("./meson.exe")
	if err != nil {
		basic.Logger.Errorln("RestartTerminal path_util.SmartExistPath err:", err)
		return err
	}
	if !exist {
		basic.Logger.Errorln("RestartNode path_util.SmartPathExist file not exist")
		return err
	}

	cmd := exec.Command("cmd", "/C", fmt.Sprintf("%s service restart", absPath))
	err = cmd.Run()
	if err != nil {
		basic.Logger.Errorln("restart peer node error:", err)
		return err
	}
	return nil
}
