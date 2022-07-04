package client

import (
	"errors"

	"github.com/meson-network/peer-node/basic/conf"
	api2 "github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/api"
)

var EndPoint string
var Token string

func Init() error {
	toml_conf := conf.Get_config().Toml_config
	token_str := toml_conf.Token
	if token_str == "" {
		return errors.New("config error : [token]")
	}

	endpoint := toml_conf.EndPoint
	if endpoint == "" {
		return errors.New("config error : [end_point]")
	}

	url := endpoint + "/api/user/token_check"
	res := &api.API_META_STATUS{}
	err := api2.Get_(url, token_str, 30, res)
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
