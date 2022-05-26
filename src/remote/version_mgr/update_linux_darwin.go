//go:build linux || darwin
// +build linux darwin

package version_mgr

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

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
	basic.Logger.Infoln("New version detected, start to upgrade... ")

	//new version download url

	//check arch and os
	arch, osInfo := GetOSInfo()

	// 'https://meson.network/static/terminal/v0.1.2/meson-darwin-amd64.tar.gz'
	fileName := "meson" + "-" + osInfo + "-" + arch + ".tar.gz"
	downloadPath := "/v" + latestVersion + "/" + fileName
	newVersionDownloadUrl := "xxxx domain" + downloadPath
	basic.Logger.Debugln("new version download url", "url", newVersionDownloadUrl)

}

func DownloadNewVersion(fileName string, downloadUrl string, newVersion string) error {
	//get
	response, err := http.Get(downloadUrl)
	if err != nil {
		//logger.Error("get file url "+downloadUrl+" error", "err", err)
		return err
	}
	//creat folder and file
	distDir := filepath.Dir(fileName)
	err = os.MkdirAll(distDir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	//defer file.Close()
	if response.Body == nil {
		os.Remove(fileName)
		file.Close()
		return errors.New("body is null")
	}
	defer response.Body.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		os.Remove(fileName)
		file.Close()
		return err
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		os.Remove(fileName)
		file.Close()
		return err
	}
	size := fileInfo.Size()
	//logger.Debug("donwload file,fileInfo", "size", size)

	if size == 0 {
		os.Remove(fileName)
		file.Close()
		return errors.New("download file size error")
	}
	file.Close()

	////unzip tar.gz
	//targetDir := strings.Replace(fileName, ".tar.gz", "", 1)
	//// file read
	//fr, err := os.Open(fileName)
	//if err != nil {
	//	logger.Error("open version file error", "err", err)
	//	return err
	//}
	//defer fr.Close()
	//// gzip read
	//gr, err := gzip.NewReader(fr)
	//if err != nil {
	//	logger.Error("gzip read new version file error", "err", err)
	//	return err
	//}
	//defer gr.Close()
	//// tar read
	//tr := tar.NewReader(gr)
	//// 读取文件
	//for {
	//	h, err := tr.Next()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		logger.Error("unzip new version file error", "err", err)
	//		return err
	//	}
	//	fileName := runpath.RunPath + "/" + h.Name
	//	dirName := string([]rune(fileName)[0:strings.LastIndex(fileName, "/")])
	//	err = os.MkdirAll(dirName, 0777)
	//	if err != nil {
	//		logger.Error("unzip new version file error-create dir", "err", err)
	//		return err
	//	}
	//	if utils.IsDir(fileName) {
	//		continue
	//	}
	//	fw, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777 /*os.FileMode(h.Mode)*/)
	//	//fw,err:=os.Create(fileName)
	//	if err != nil {
	//		logger.Error("unzip new version file error-create file", "err", err)
	//		return err
	//	}
	//	defer fw.Close()
	//	// 写文件
	//	_, err = io.Copy(fw, tr)
	//	if err != nil {
	//		logger.Error("unzip new version file error-copy file", "err", err)
	//		return err
	//	}
	//}
	//logger.Debug("un tar.gz ok")
	//
	////cover old version file
	//files, _ := ioutil.ReadDir(targetDir)
	//for _, v := range files {
	//	fileName := v.Name()
	//	fmt.Println(v.Name())
	//	if fileName == "config.txt" {
	//		continue
	//	}
	//	oldFile := path.Join(runpath.RunPath, fileName)
	//	newFile := path.Join(targetDir, fileName)
	//	err := coverOldFile(newFile, oldFile)
	//	if err != nil {
	//		logger.Error("new version file error-cover file", "err", err)
	//		continue
	//	}
	//}
	//
	//versionFile := path.Join(runpath.RunPath, "./v"+Version)
	//os.Remove(versionFile)
	//
	//os.RemoveAll(targetDir)
	//os.Remove(fileName)
	//
	//return nil

	return nil
}

func coverOldFile(newFile string, oldFile string) error {
	input, err := ioutil.ReadFile(newFile)
	if err != nil {
		return err
	}
	os.Remove(oldFile)
	err = ioutil.WriteFile(oldFile, input, 777)
	if err != nil {
		fmt.Println("Error creating", oldFile)
		fmt.Println(err)
		return err
	}
	os.Chmod(oldFile, 0777)
	return nil
}

func RestartTerminal() {
	basic.Logger.Debugln("peer node restart cmd")

	absPath, err := path_util.SmartExistPath("./peer-node")
	if err != nil {
		basic.Logger.Errorln("RestartTerminal path_util.SmartExistPath err:", err)
		return
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("sudo %s service restart", absPath))
	err = cmd.Run()
	if err != nil {
		basic.Logger.Errorln("restart peer node error:", err)
	}
}
