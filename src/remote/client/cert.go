package client

import (
	"errors"

	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/dns"
)

func GetCertHash() (string, error) {
	hashRes := &dns.Msg_Resp_CertHash{}
	err := api.Get_(EndPoint+"/api/node/cert/hash", Token, 30, hashRes)
	if err != nil {
		return "", err
	}
	if hashRes.Meta_status <= 0 {
		return "", errors.New(hashRes.Meta_message)
	}

	return hashRes.Hash, nil
}

func GetCert() (crt string, key string, err error) {
	res := &dns.Msg_Resp_Cert{}
	err = api.Get_(EndPoint+"/api/node/cert", Token, 30, res)
	if err != nil {
		return "", "", err
	}

	if res.Meta_status <= 0 {
		return "", "", errors.New(res.Meta_message)
	}

	return res.Crt, res.Key, nil
}
