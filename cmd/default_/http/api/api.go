package api

import (
	_ "github.com/meson-network/peer-node/cmd/default_/http/api_docs"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
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
