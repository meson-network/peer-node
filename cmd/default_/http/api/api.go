package api

import (
	"github.com/meson-network/peer-node/basic"

	"github.com/coreservice-io/utils/path_util"
	_ "github.com/meson-network/peer-node/cmd/default_/http/api_docs"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/gen"
)

// for swagger
// @title           api example
// @version         1.0
// @description     api example
// @termsOfService  https://domain.com
// @contact.name    Support
// @contact.url     https://domain.com
// @contact.email   contact@domain.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes         https

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

func DeclareApi(httpServer *echo_plugin.EchoServer) {
	//health
	//never change this api
	//there will be callback from checker
	//there won't be any token reward if fails
	//there maybe credit punishment if node health fails(not active)
	config_health(httpServer)

	config_download(httpServer)
	config_heart_beat_callback(httpServer)
	config_speed_tester(httpServer)
}

func ConfigApi(httpServer *echo_plugin.EchoServer) {
	httpServer.GET("/swagger/*", echoSwagger.WrapHandler)
}

func Gen_Api_Docs() {

	api_doc_gen_search_dir, _ := configuration.Config.GetString("api_doc_gen_search_dir", "")
	api_doc_gen_mainfile, _ := configuration.Config.GetString("api_doc_gen_mainfile", "")
	api_doc_gen_output_dir, _ := configuration.Config.GetString("api_doc_gen_output_dir", "")

	if api_doc_gen_search_dir == "" ||
		api_doc_gen_mainfile == "" ||
		api_doc_gen_output_dir == "" {
		basic.Logger.Errorln("api_doc_gen_search_dir|api_doc_gen_mainfile|api_doc_gen_output_dir config errors")
		return
	}

	api_f, api_f_err := path_util.SmartExistPath(api_doc_gen_search_dir)
	if api_f_err != nil {
		basic.Logger.Errorln("api_doc_gen_search_dir folder not exist")
		return
	}
	api_doc_f, api_doc_f_err := path_util.SmartExistPath(api_doc_gen_output_dir)
	if api_doc_f_err != nil {
		basic.Logger.Errorln("api_doc_gen_output_dir folder not exist")
		return
	}

	config := &gen.Config{
		SearchDir:       api_f,
		OutputDir:       api_doc_f,
		MainAPIFile:     api_doc_gen_mainfile,
		OutputTypes:     []string{"go", "json", "yaml"},
		ParseDependency: true,
	}

	err := gen.New().Build(config)
	if err != nil {
		basic.Logger.Errorln("Gen_Api_Docs", err)
	}

}
