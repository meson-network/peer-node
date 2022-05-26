//go:build windows
// +build windows

package version_mgr

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

	// 'https://meson.network/static/terminal/v0.1.2/meson-windows-amd64.zip'
	fileName := "meson" + "-" + osInfo + "-" + arch + ".zip"
	downloadPath := "/v" + latestVersion + "/" + fileName
	newVersionDownloadUrl := "xxxx domain" + downloadPath
	basic.Logger.Debugln("new version download url", "url", newVersionDownloadUrl)

}

func RestartTerminal() {
	basic.Logger.Debugln("peer node restart cmd")

	absPath, err := path_util.SmartExistPath("./peer-node.exe")
	if err != nil {
		basic.Logger.Errorln("RestartTerminal path_util.SmartExistPath err:", err)
		return
	}

	cmd := exec.Command("cmd", "/C", fmt.Sprintf("%s service restart", absPath))
	err = cmd.Run()
	if err != nil {
		basic.Logger.Errorln("restart peer node error:", err)
	}
}
