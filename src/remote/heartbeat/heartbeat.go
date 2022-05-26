package heartbeat

import (
	"errors"

	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/heart_beat"
)

var server_unix_time int64

func GetLastHeartTime() int64 {
	return server_unix_time
}

func SetLastHeartTime(time int64) {
	server_unix_time = time
}

func SendHeartBeat(hb_req *heart_beat.Msg_Req_HeartBeat) (*heart_beat.Msg_Resp_HeartBeat, error) {
	res := &heart_beat.Msg_Resp_HeartBeat{}
	err := api.POST_(client.EndPoint+"/api/node/heartbeat", client.Token, hb_req, 30, res)
	if err != nil {
		return nil, err
	}

	if res.Meta_status <= 0 {
		return nil, errors.New(res.Meta_message)
	}
	SetLastHeartTime(res.Server_unixtime)
	return res, nil
}
