package default_

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/cmd/default_/plugin"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/access_key_mgr"
	"github.com/meson-network/peer-node/src/cdn_cache_folder"
	"github.com/meson-network/peer-node/src/cert_mgr"
	"github.com/meson-network/peer-node/src/common/dbkv"
	"github.com/meson-network/peer-node/src/file_mgr"
	"github.com/meson-network/peer-node/src/node_info"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/src/schedule_job"
	"github.com/meson-network/peer-node/src/speed_tester_file"
	"github.com/meson-network/peer-node/src/version_mgr"
	"github.com/urfave/cli/v2"
)

func StartDefault(clictx *cli.Context) {

	//err := cert_mgr.Init()
	//if err != nil {
	//	basic.Logger.Fatalln("initCert error", err)
	//}
	//RunMinio()
	//return

	color.Green(basic.Logo)
	color.Green(fmt.Sprintf("Node Version: v%s", version_mgr.NodeVersion))

	//init cdn cache folder
	err := cdn_cache_folder.Init()
	if err != nil {
		basic.Logger.Fatalln("init cdn cache folder err:", err)
	}

	///////////////////
	plugin.InitPlugin()
	///////////////////

	//token check first
	err = client.Init()
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	//version_mgr
	err = version_mgr.Init()
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	//accessKey
	curKey := ""
	preKey := ""
	accKey, err := dbkv.GetKey(sqlite_plugin.GetInstance(), "access_key", false, false)
	if err == nil && accKey != "" {
		keys := strings.Split(accKey, ",")
		if len(keys) == 2 {
			curKey = keys[0]
			preKey = keys[1]
		}
	}
	err = access_key_mgr.Init(curKey, preKey)
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	//clean not finished download job and files
	err = file_mgr.CleanDownloadingFiles()
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	//check cache folder
	err = cdn_cache_folder.GetInstance().CheckFolder(1)
	if err != nil {
		basic.Logger.Fatalln("check cdn cache folder err:", err)
	}
	speed_tester_file.CheckTesterFile()

	//init node
	err = node_info.InitNode()
	if err != nil {
		basic.Logger.Fatalln("initNode error", err)
	}

	////////init update cert
	err = cert_mgr.Init()
	if err != nil {
		basic.Logger.Fatalln("initCert error", err)
	}

	cert_update_err := cert_mgr.GetInstance().UpdateCert(nil)
	if cert_update_err != nil {
		basic.Logger.Fatalln("init certificate error", cert_update_err)
	}
	///////////////////////////////

	//init httpserver
	err_server := plugin.InitEchoServer()
	if err_server != nil {
		basic.Logger.Fatalln(err_server)
	}

	//start the httpserver
	go http.StartDefaultHttpSever()

	//start jobs
	go start_jobs()

	for {
		//never quit
		time.Sleep(time.Duration(1) * time.Hour)
	}

}

func start_jobs() {
	//start threads jobs
	//check all services already started
	if !http.CheckDefaultHttpServerStarted() {
		basic.Logger.Fatalln("http server not working")
	}

	/////////
	schedule_job.CheckVersion()
	schedule_job.ScanExpirationFile()
	schedule_job.UpdateCert()
	schedule_job.ScanLeakFile()
	schedule_job.RenewAccessKey()
	schedule_job.DeleteEmptyFolder()

	schedule_job.HeartBeat()
}
