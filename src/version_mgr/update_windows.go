//go:build windows
// +build windows

package version_mgr

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
)

func genFileName() string {
	//check arch and os
	arch, osInfo := GetOSInfo()
	// 'https://xxxx/xxxxx/node/v0.1.2/meson-windows-amd64.zip'
	return "meson" + "-" + osInfo + "-" + arch + ".zip"
}

func unzip(targetFolder string, body io.Reader) error {
	//get content
	tempContent, err := ioutil.ReadAll(body)
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

		filePath := filepath.Join(targetFolder, name)
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

		err = ioutil.WriteFile(filePath, content, 777)
		if err != nil {
			basic.Logger.Errorln("Error creating", filePath, "err:", err)
			return err
		}
		os.Chmod(filePath, 0777)
	}
	basic.Logger.Debugln("un zip ok")
	return nil
}

func overwriteOldFile(newFile string, oldFile string) error {
	input, err := ioutil.ReadFile(newFile)
	if err != nil {
		return err
	}
	//rename oldFile
	//todo how to remove opened file
	reNameFile := oldFile + ".old"
	os.Rename(oldFile, reNameFile)
	err = ioutil.WriteFile(oldFile, input, 777)
	if err != nil {
		basic.Logger.Errorln("windows overwriteOldFile WriteFile error:", err, "filepath:", oldFile)
		os.Rename(reNameFile, oldFile)
		return err
	}
	os.Remove(reNameFile)
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
