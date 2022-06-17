package precheck_config

import (
	"fmt"
	"strconv"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	api2 "github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/api"
	"github.com/meson-network/peer_common/cdn_cache"
)

func PreCheckConfig() {
	toml_conf := conf.Get_config().Toml_config

	token := toml_conf.Token
	if len(token) != 24 {
		basic.Logger.Fatalln("token config error")
	}

	port := toml_conf.Https_port
	if port <= 0 || port > 65535 || echo_plugin.IsForbiddenPort(port) {
		basic.Logger.Fatalln("https_port cofig error")
	}

	cdnCacheSize := toml_conf.Cache.Size
	if cdnCacheSize < cdn_cache.MIN_CACHE_SIZE {
		basic.Logger.Fatalln("cache.size config error,minimum is 20")
	}

}

func CheckConfig() {
	checkToken()
	checkPort()
	checkProvideCacheSize()

	conf.Get_config().Save_config()
}

func checkToken() {
	toml_conf := conf.Get_config().Toml_config

	endpoint := toml_conf.EndPoint
	if endpoint == "" {
		basic.Logger.Fatalln("[end_point] in config error")
		return
	}

	token := toml_conf.Token
	needInputToken := false
	for {
		if needInputToken {
			var myToken string
			fmt.Println("can not find your token. Please login https://meson.network")
			fmt.Printf("Please enter your token: ")
			_, err := fmt.Scanln(&myToken)
			if err != nil {
				fmt.Println("read input token error")
			}
			token = myToken
		}
		if len(token) != 24 {
			needInputToken = true
			fmt.Println("token length error")
			continue
		}

		//check token
		url := endpoint + "/api/user/token_check"
		res := &api.API_META_STATUS{}
		err := api2.Get(url, token, res)
		if err != nil {
			fmt.Println("check token error:", err)
			needInputToken = true
			continue
		}
		if res.Meta_status <= 0 {
			fmt.Println(res.Meta_message)
			needInputToken = true
			continue
		}
		fmt.Println("token:", token)

		toml_conf.Token = token
		toml_conf.Storage.Password = token
		break
	}
}

func checkPort() {

	toml_conf := conf.Get_config().Toml_config
	port := toml_conf.Https_port

	for {
		if port <= 0 || port > 65535 || echo_plugin.IsForbiddenPort(port) {
			//input port
			fmt.Printf("CAN NOT use port %d ,please input port, \n", port)
		} else {
			fmt.Println("port:", port)
			toml_conf.Https_port = port
			break
		}
		var myPortStr string
		fmt.Printf("please input port(default is 443):")
		fmt.Scanln(&myPortStr)
		if myPortStr == "" {
			port = 443
		} else {
			var err error
			port, err = strconv.Atoi(myPortStr)
			if err != nil {
				fmt.Println("input error:", err)
			}
		}
	}
}

func checkProvideCacheSize() {
	toml_conf := conf.Get_config().Toml_config
	cdnCacheSize := toml_conf.Cache.Size

	for {
		if cdnCacheSize < cdn_cache.MIN_CACHE_SIZE {
			//input port
			fmt.Printf("cache_size must be at least %d GB \n", cdn_cache.MIN_CACHE_SIZE)
		} else {
			fmt.Println("cache_size:", cdnCacheSize)
			toml_conf.Cache.Size = cdnCacheSize
			break
		}
		var sizeStr string
		fmt.Printf("please input provide disk size GB(at least %d,default is %d):", cdn_cache.MIN_CACHE_SIZE, cdn_cache.MIN_CACHE_SIZE)
		fmt.Scanln(&sizeStr)
		if sizeStr == "" {
			cdnCacheSize = cdn_cache.MIN_CACHE_SIZE
		} else {
			var err error
			cdnCacheSize, err = strconv.Atoi(sizeStr)
			if err != nil {
				fmt.Println("input error:", err)
			}
		}
	}
}
