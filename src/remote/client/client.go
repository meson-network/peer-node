package client

import (
	"errors"

	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/tools/http"
	"github.com/meson-network/peer_common/api"
)

type Client struct {
	EndPoint string
	Token    string
}

var client *Client

func GetClient() (*Client, error) {
	if client != nil {
		return client, nil
	}

	new_client, new_client_err := new_client()
	if new_client_err != nil {
		return nil, new_client_err
	}

	client = new_client
	return client, nil
}

func new_client() (*Client, error) {

	token_str, err := configuration.Config.GetString("token", "")
	if err != nil || token_str == "" {
		return nil, errors.New("config error : token [string] in config.json ")
	}

	endpoint := "https://server.mesontracking.com"
	url := endpoint + "/api/token/check"
	res := &api.API_META_STATUS{}
	error := http.Get(url, token_str, res)
	if error != nil {
		return nil, error
	}

	if res.Meta_status <= 0 {
		return nil, errors.New(res.Meta_message)
	}

	client := &Client{
		EndPoint: endpoint,
		Token:    token_str,
	}

	return client, nil
}
