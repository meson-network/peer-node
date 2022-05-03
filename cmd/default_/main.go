package default_

import (
	"github.com/fatih/color"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/cmd/default_/plugin"
	"github.com/meson-network/peer-node/src/downloader_mgr"
	"github.com/meson-network/peer-node/src/storage_mgr"
	"github.com/urfave/cli/v2"
)

func StartDefault(clictx *cli.Context) {

	//defer func() {
	//	//global.ReleaseResources()
	//}()
	color.Green(basic.Logo)

	ini_components()

	//start threads jobs
	go start_jobs()

	start()
}

func start() {
	//start the httpserver
	http.StartDefaultHttpSever()
}

func ini_components() {
	//ini components and run example
	plugin.InitPlugin()

	//storagemgr
	err := storage_mgr.Init()
	if err != nil {
		panic(err)
	}
}

func start_jobs() {
	//check all services already started
	if !http.CheckDefaultHttpServerStarted() {
		panic("http server not working")
	}

	//test downloader
	downloader_mgr.StartDownloader("https://dl.google.com/go/go1.18.1.darwin-amd64.pkg", func(filehash, file_path string) {
		basic.Logger.Infoln("download_success filehash:", filehash)
		basic.Logger.Infoln("download_success file_path:", file_path)
	}, func(filehash string, download_code int) {
		basic.Logger.Infoln("download_failure filehash:", filehash)
		basic.Logger.Infoln("download_failure download_code:", download_code)
	})

}
