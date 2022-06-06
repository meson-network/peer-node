package client

import (
	"errors"

	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/node_domain"
)

func GetNodeDomain() (string, error) {
	result := &node_domain.Msg_Resp_NodeDomain{}
	err := api.Get_(EndPoint+"/api/node/node_domain", Token, 30, result)
	if err != nil {
		return "", err
	}

	if result.Meta_status <= 0 {
		return "", errors.New(result.Meta_message)
	}

	return result.Node_domain, nil
}
