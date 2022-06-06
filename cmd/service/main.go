package service

import (
	"fmt"
	"os"
	"strconv"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/plugin/daemon_plugin"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	api2 "github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/api"
	"github.com/meson-network/peer_common/cdn_cache"
	"github.com/urfave/cli/v2"
)

func RunServiceCmd(clictx *cli.Context) {

	daemon_name, err := configuration.Config.GetString("daemon_name", "meson-node")
	if err != nil {
		basic.Logger.Errorln("daemon_name [string] in config error," + err.Error())
		return
	}

	if daemon_name == "" {
		basic.Logger.Errorln("daemon_name in config should not be vacant")
		return
	}

	exe_path, exe_path_err := os.Executable()
	if exe_path_err != nil {
		basic.Logger.Errorln(exe_path_err)
		return

	}

	//exeDir := filepath.Dir(exe_path)
	//
	//if _, dir_err := os.Stat(path.Join(exeDir, "assets")); dir_err != nil {
	//	basic.Logger.Errorln("error -> please check:")
	//	basic.Logger.Errorln("1.dont directly `go run` for service, always `go build` first")
	//	basic.Logger.Errorln("2.the assets folder exist parellel to the excutable file ")
	//	return
	//}

	basic.Logger.Infoln("exefile:" + exe_path + " to be service target")

	//check command
	subCmds := clictx.Command.Names()
	if len(subCmds) == 0 {
		basic.Logger.Fatalln("no sub command")
	}

	action := subCmds[0]
	err = daemon_plugin.Init(daemon_name)
	if err != nil {
		basic.Logger.Fatalln("init daemon service error:", err)
	}

	var status string
	var e error
	switch action {
	case "install":
		//check config
		checkConfig()
		status, e = daemon_plugin.GetInstance(daemon_name).Install()
		basic.Logger.Debugln("cmd install")
	case "remove":
		daemon_plugin.GetInstance(daemon_name).Stop()
		status, e = daemon_plugin.GetInstance(daemon_name).Remove()
		basic.Logger.Debugln("cmd remove")
	case "start":
		//check config
		checkConfig()
		status, e = daemon_plugin.GetInstance(daemon_name).Start()
		basic.Logger.Debugln("cmd start")
	case "stop":
		status, e = daemon_plugin.GetInstance(daemon_name).Stop()
		basic.Logger.Debugln("cmd stop")
	case "restart":
		daemon_plugin.GetInstance(daemon_name).Stop()
		//check config
		checkConfig()
		status, e = daemon_plugin.GetInstance(daemon_name).Start()
		basic.Logger.Debugln("cmd restart")
	case "status":
		status, e = daemon_plugin.GetInstance(daemon_name).Status()
		basic.Logger.Debugln("cmd status")
	default:
		basic.Logger.Debugln("no sub command")
		return
	}

	if e != nil {
		fmt.Println(status, "\nError: ", e)
		os.Exit(1)
	}
	fmt.Println(status)
}

//{
//"token":"mytoken",
//"cdn_cache_size":30,
//"https_port": 443
//}

func checkConfig() {
	checkToken()
	checkPort()
	checkProvideCacheSize()

	configuration.Config.WriteConfig()
}

func checkToken() {
	endpoint, err := configuration.Config.GetString("endpoint", "https://api.meson.network")
	if err != nil || endpoint == "" {
		basic.Logger.Errorln("endpoint [string] in config error," + err.Error())
		return
	}

	token, err := configuration.Config.GetString("token", "")
	if err != nil {
		basic.Logger.Errorln("daemon_name [string] in config error," + err.Error())
		return
	}

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
		err = api2.Get(url, token, res)
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
		configuration.Config.Set("token", token)
		break
	}
}

func checkPort() {
	port, err := configuration.Config.GetInt("https_port", 0)
	if err != nil {
		basic.Logger.Errorln("https_port [int] in config error," + err.Error())
		return
	}

	for {
		if port <= 0 || port > 65535 || echo_plugin.IsForbiddenPort(port) {
			//input port
			fmt.Printf("CAN NOT use port %d ,please input port, \n", port)
		} else {
			fmt.Println("port:", port)
			configuration.Config.Set("https_port", port)
			break
		}
		var myPortStr string
		fmt.Printf("please input port(default is 443):")
		fmt.Scanln(&myPortStr)
		if myPortStr == "" {
			port = 443
		} else {
			port, err = strconv.Atoi(myPortStr)
		}
	}
}

func checkProvideCacheSize() {
	cdnCacheSize, err := configuration.Config.GetInt("cache_size", 0)
	if err != nil {
		basic.Logger.Errorln("cdn_cache_size [string] in config error," + err.Error())
		return
	}

	for {
		if cdnCacheSize < cdn_cache.MIN_CACHE_SIZE {
			//input port
			fmt.Printf("cache_size must be at least %d GB \n", cdn_cache.MIN_CACHE_SIZE)
		} else {
			fmt.Println("cache_size:", cdnCacheSize)
			configuration.Config.Set("cache_size", cdnCacheSize)
			break
		}
		var sizeStr string
		fmt.Printf("please input provide disk size GB(at least %d,default is %d):", cdn_cache.MIN_CACHE_SIZE, cdn_cache.MIN_CACHE_SIZE)
		fmt.Scanln(&sizeStr)
		if sizeStr == "" {
			cdnCacheSize = cdn_cache.MIN_CACHE_SIZE
		} else {
			cdnCacheSize, err = strconv.Atoi(sizeStr)
		}
	}
}
