package default_

import (
	"time"

	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/cmd/default_/plugin"
	"github.com/meson-network/peer-node/src/remote/cert"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/src/storage_mgr"
	"github.com/urfave/cli/v2"
)

func StartDefault(clictx *cli.Context) {

	color.Green(basic.Logo)

	plugin.InitPlugin()

	//token check first
	_, c_err := client.GetClient()
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

	//init storage
	stor_err := storage_mgr.Init()
	if stor_err != nil {
		basic.Logger.Fatalln(stor_err)
	}
	///////////////////

	/////////////////////////
	err_server := plugin.InitEchoServer()
	if err_server != nil {
		basic.Logger.Fatalln(err_server)
	}

	//////////////start the httpserver
	go http.StartDefaultHttpSever()

	/////////////////start jobs
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
	cert_m, cert_m_err := cert.GetCertMgr()
	if cert_m_err != nil {
		basic.Logger.Fatalln(cert_m_err)
	}
	cert_m.ScheduleUpdateJob(func(crt, key string) {
		http.ServerReloadCert()
	})
}
