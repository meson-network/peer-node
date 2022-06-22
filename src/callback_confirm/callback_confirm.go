package callback_confirm

import (
	"strconv"
	"time"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/remote/client"
)

var HeartBeatReceiveChan = make(chan bool)
var WaitingHeartBeatCallback = true

func WaitHeartBeatCallbackConfirm() {
	go func() {
		basic.Logger.Infoln("waiting for heart beat callback...")
		select {
		case <-HeartBeatReceiveChan:
			basic.Logger.Infoln("heart beat callback received")
		case <-time.After(75 * time.Second):
			basic.Logger.Errorln("Heart beat callback not received. Please confirm that your machine can be accessed by the external network and the port is opened on the firewall.")
			//get domain from remote
			nodeDomain, err := client.GetNodeDomain()
			if err != nil {
				basic.Logger.Errorln("get node domain error," + err.Error())
				return
			}
			portStr := echo_plugin.GetInstance().Http_port
			checkUrl := "https://" + nodeDomain + ":" + strconv.Itoa(portStr) + "/api/health"
			basic.Logger.Infoln("check accessible by url:", checkUrl)
		}
		WaitingHeartBeatCallback = false
		close(HeartBeatReceiveChan)
	}()

}
