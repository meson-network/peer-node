package client

import (
	"errors"

	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/heart_beat"
)

func SendHeartBeat(hb_req *heart_beat.Msg_Req_HeartBeat) (*heart_beat.Msg_Resp_HeartBeat, error) {
	res := &heart_beat.Msg_Resp_HeartBeat{}
	err := api.POST_(EndPoint+"/api/node/heartbeat", Token, hb_req, 30, res)
	if err != nil {
		return nil, err
	}

	if res.Meta_status <= 0 {
		return nil, errors.New(res.Meta_message)
	}

	return res, nil
}
