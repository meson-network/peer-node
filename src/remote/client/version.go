package client

import (
	"errors"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/version"
)

func GetNodeVersionFromServer() (latestVersion string, allowVersion string, err error) {
	//check is there new version or not
	basic.Logger.Debugln("Check Version...")
	result := &version.Msg_Resp_NodeVersion{}
	err = api.Get_(EndPoint+"/api/node/version", Token, 30, result)
	if err != nil {
		return "", "", err
	}

	if result.Meta_status <= 0 {
		return "", "", errors.New(result.Meta_message)
	}

	return result.Latest_version, result.Allow_version, nil
}
