package default_

import (
	"time"

	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/cmd/default_/plugin"
	"github.com/meson-network/peer-node/src/api/cert"
	"github.com/meson-network/peer-node/src/api/client"
	"github.com/urfave/cli/v2"
)

func StartDefault(clictx *cli.Context) {

	color.Green(basic.Logo)

	plugin.InitPlugin()

	//token check first
	cient_, c_err := client.GetClient()
	if c_err != nil {
		basic.Logger.Fatalln(c_err)
	}

	//update cert
	cert_m, cert_m_err := cert.GetCertMgr()
	if cert_m_err != nil {
		basic.Logger.Fatalln(cert_m_err)
	}

	cert_update_err := cert_m.UpdateCert(cient_)
	if cert_update_err != nil {
		basic.Logger.Fatalln(cert_update_err)
	}

	// err := storage_mgr.Init()
	// if err != nil {
	// 	panic(err)
	// }

	/////////////////////////
	err := plugin.InitEchoServer()
	if err != nil {
		basic.Logger.Fatalln(err)
	}

	//get cert from remote

	//start the httpserver
	go http.StartDefaultHttpSever()

	//start threads jobs
	//check all services already started
	if !http.CheckDefaultHttpServerStarted() {
		panic("http server not working")
	}

	for {
		//never quit
		time.Sleep(time.Duration(1) * time.Hour)
	}

}
