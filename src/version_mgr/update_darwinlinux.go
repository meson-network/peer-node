//go:build linux || darwin
// +build linux darwin

package version_mgr

import (
	"archive/tar"
	"compress/gzip"
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
	// 'https://xxx.xxxx/xxxx/node/v0.1.2/meson-darwin-amd64.tar.gz'
	return "meson" + "-" + osInfo + "-" + arch + ".tar.gz"
}

func unzip(targetFolder string, body io.Reader) error {
	// gzip read
	gr, err := gzip.NewReader(body)
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

		filePath := filepath.Join(targetFolder, name)
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

		err = ioutil.WriteFile(filePath, content, 0777)
		if err != nil {
			basic.Logger.Errorln("Error creating", filePath, "err:", err)
			return err
		}
	}
	basic.Logger.Debugln("un tar.gz ok")
	return nil
}

func overwriteOldFile(newFile string, oldFile string) error {
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
	return nil
}

func RestartNode() error {
	basic.Logger.Debugln("peer node restart cmd")

	absPath, exist, err := path_util.SmartPathExist("./meson")
	if err != nil {
		basic.Logger.Errorln("RestartNode path_util.SmartPathExist err:", err)
		return err
	}
	if !exist {
		basic.Logger.Errorln("RestartNode path_util.SmartPathExist file not exist")
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
