package node_config

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/http/api"
	common_node_config "github.com/meson-network/peer_common/node_config"
)

var (
	Cached_file_expire_secs = int64(3600 * 6)
	Max_file_size_bytes     = int64(1024 * 1024 * 1024)
	Free_space_line         = int64(1024 * 1024 * 1024)
	Delete_trigger_rate     = float64(0.7)
)

func GetNodeConfig() {
	res := &common_node_config.Msg_Resp_NodeConfig{}
	err := api.Get_(client.EndPoint+"/api/node/node_config", client.Token, 30, res)
	if err != nil {
		basic.Logger.Errorln("reportExpiredFiles post error:", err)
		return
	}

	if res.Meta_status <= 0 {
		basic.Logger.Errorln("reportExpiredFiles post error:", res.Meta_message)
		return
	}

	if res.Cached_file_expire_secs > 0 {
		Cached_file_expire_secs = res.Cached_file_expire_secs
	}
	if res.Max_file_size_bytes > 0 {
		Max_file_size_bytes = res.Max_file_size_bytes
	}
	if res.Free_space_line > 0 {
		Free_space_line = res.Free_space_line
	}
	if res.Delete_trigger_rate > 0 {
		Delete_trigger_rate = res.Delete_trigger_rate
	}
}
