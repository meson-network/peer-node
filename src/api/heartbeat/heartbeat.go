package heartbeat

import (
	"errors"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/api/client"
	"github.com/meson-network/peer-node/tools/http"
	"github.com/meson-network/peer_common/heart_beat"
)

var server_unix_time int64

func GetLastHeartTime() int64 {
	return server_unix_time
}

func SetLastHeartTime(time int64) {
	server_unix_time = time
}

func RequestHeartBeat(hb_req *heart_beat.Msg_Req_HeartBeat) (*heart_beat.Msg_Resp_HeartBeat, error) {

	cient_, c_err := client.GetClient()
	if c_err != nil {
		basic.Logger.Fatalln(c_err)
	}

	res := &heart_beat.Msg_Resp_HeartBeat{}
	err := http.POST(cient_.EndPoint+"/api/node/heartbeat", cient_.Token, hb_req, res)
	if err != nil {
		return nil, err
	}

	if res.Meta_status <= 0 {
		return nil, errors.New(res.Meta_message)
	}
	SetLastHeartTime(res.Server_unixtime)
	return res, nil
}
