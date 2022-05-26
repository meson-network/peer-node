package default_

import (
	"time"

	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/cmd/default_/plugin"
	"github.com/meson-network/peer-node/src/cdn_cache_folder"
	"github.com/meson-network/peer-node/src/info"
	"github.com/meson-network/peer-node/src/remote/cert"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/urfave/cli/v2"
)

func StartDefault(clictx *cli.Context) {

	//RunMinio()
	//
	//return

	color.Green(basic.Logo)

	//init cdn cache folder
	err := cdn_cache_folder.Init()
	if err != nil {
		basic.Logger.Fatalln("init cdn cache folder err:", err)
	}

	///////////////////
	plugin.InitPlugin()
	///////////////////

	//delete

	err = cdn_cache_folder.GetInstance().CheckFolder(10)
	if err != nil {
		basic.Logger.Fatalln("check cdn cache folder err:", err)
	}

	//clean not finished download job and files

	//scan db record clean files which not exist on disk

	//scan folder clean file which not in db

	//init node
	err = info.InitNode()
	if err != nil {
		basic.Logger.Fatalln("initNode error", err)
	}

	//token check first
	c_err := client.Init()
	if c_err != nil {
		basic.Logger.Fatalln(c_err)
	}

	////////init update cert
	cert_m, cert_m_err := cert.GetCertMgr()
	if cert_m_err != nil {
		basic.Logger.Fatalln(cert_m_err)
	}

	cert_update_err := cert_m.UpdateCert(nil)
	if cert_update_err != nil {
		basic.Logger.Fatalln(cert_update_err)
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
	//todo get cert again???
	cert_m, cert_m_err := cert.GetCertMgr()
	if cert_m_err != nil {
		basic.Logger.Fatalln(cert_m_err)
	}
	cert_m.ScheduleUpdateJob(func(crt, key string) {
		http.ServerReloadCert()
	})

	//test a download task
	//download_mgr.DoTask(func(filehash string, file_local_abs_path string) {
	//	basic.Logger.Infoln("sucess download task callback", filehash, file_local_abs_path)
	//}, func(filehash string, download_code int) {
	//	basic.Logger.Infoln("failed download task callback", filehash, download_code)
	//})
}
