package client

import (
	"errors"

	"github.com/meson-network/peer-node/configuration"
	api2 "github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/api"
)

var EndPoint string
var Token string

func Init() error {
	token_str, err := configuration.Config.GetString("token", "")
	if err != nil || token_str == "" {
		return errors.New("config error : token [string] in config.json ")
	}

	endpoint := "https://server.mesontracking.com"
	url := endpoint + "/api/user/token_check"
	res := &api.API_META_STATUS{}
	err = api2.Get(url, token_str, res)
	if err != nil {
		return err
	}

	if res.Meta_status <= 0 {
		return errors.New(res.Meta_message)
	}

	EndPoint = endpoint
	Token = token_str

	return nil
}
