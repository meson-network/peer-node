package cert

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/src/api/client"
	"github.com/meson-network/peer-node/tools/http"
	"github.com/meson-network/peer_common/dns"
)

type CertMgr struct {
	Key_path string
	Crt_path string
}

var cert_mgr *CertMgr

func GetCertMgr() (*CertMgr, error) {
	if cert_mgr != nil {
		return cert_mgr, nil
	}

	rel_crt, err := configuration.Config.GetString("https_crt_path", "")
	if err != nil || rel_crt == "" {
		return nil, errors.New("https_crt_path [string] in config.json err")
	}

	rel_key, err := configuration.Config.GetString("https_key_path", "")
	if err != nil || rel_key == "" {
		return nil, errors.New("https_key_path [string] in config.json err")
	}

	crt_path, cert_path_err := path_util.SmartExistPath(rel_crt)
	if cert_path_err != nil {
		return nil, errors.New("https crt file path error," + cert_path_err.Error())
	}

	key_path, key_path_err := path_util.SmartExistPath(rel_key)
	if cert_path_err != nil {
		return nil, errors.New("https key file path error," + key_path_err.Error())
	}

	cert_mgr = &CertMgr{
		key_path,
		crt_path,
	}
	return cert_mgr, nil
}

func (c *CertMgr) UpdateCert(client *client.Client) error {

	res := &dns.Msg_Resp_Cert{}
	err := http.Get(client.EndPoint+"/api/node/cert", client.Token, res)

	if err != nil {
		return err
	}

	if res.Meta_status <= 0 {
		return errors.New(res.Meta_message)
	}

	basic.Logger.Infoln(res)

	///////////////
	change := false
	old_crt_content, read_err := ioutil.ReadFile(c.Crt_path)
	if read_err != nil {
		return read_err
	} else {
		if string(old_crt_content) != res.Crt {
			change = true
		}
	}

	//read old .key
	old_key_content, read_err := ioutil.ReadFile(c.Key_path)
	if read_err != nil {
		return read_err
	} else {
		if string(old_key_content) != res.Key {
			change = true
		}
	}

	//update the file
	if change {
		crt_file_err := file_overwrite(c.Crt_path, res.Crt)
		if crt_file_err != nil {
			return crt_file_err
		}

		key_file_err := file_overwrite(c.Key_path, res.Key)
		if key_file_err != nil {
			return key_file_err
		}
	}

	return nil
}

func file_overwrite(path string, content string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, werr := f.WriteString(content)
	if werr != nil {
		return werr
	}
	return nil
}
