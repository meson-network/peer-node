package cert_mgr

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreservice-io/utils/hash_util"
	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/file"
)

type CertMgr struct {
	Key_path string
	Crt_path string
}

var cert_mgr *CertMgr

func Init() error {

	//todo if cert file not exist, create folder
	rel_crt, err := configuration.Config.GetString("https_crt_path", "assets/cert/public.crt")
	if err != nil || rel_crt == "" {
		return errors.New("https_crt_path [string] in config.json err")
	}

	rel_key, err := configuration.Config.GetString("https_key_path", "assets/cert/private.key")
	if err != nil || rel_key == "" {
		return errors.New("https_key_path [string] in config.json err")
	}

	crtPath := path_util.ExE_Path(rel_crt)
	crtFolder := filepath.Dir(crtPath)
	err = os.MkdirAll(crtFolder, 0777)
	if err != nil {
		return err
	}

	keyPath := path_util.ExE_Path(rel_key)
	keyFolder := filepath.Dir(keyPath)
	err = os.MkdirAll(keyFolder, 0777)
	if err != nil {
		return err
	}

	cert_mgr = &CertMgr{
		Crt_path: crtPath,
		Key_path: keyPath,
	}

	return nil
}

func GetInstance() *CertMgr {
	return cert_mgr
}

//success_callback func(string crt, string key)
func (c *CertMgr) UpdateCert(success_callback func(string, string)) error {

	//check hash
	certHash, err := client.GetCertHash()
	if err != nil {
		return err
	}

	old_crt_content, read_err := ioutil.ReadFile(c.Crt_path)
	if read_err != nil {
		return read_err
	}
	old_key_content, read_err := ioutil.ReadFile(c.Key_path)
	if read_err != nil {
		return read_err
	}

	old_content_hash := hash_util.MD5HashString(string(old_crt_content) + string(old_key_content))
	if old_content_hash == certHash {
		return nil
	}

	//need update
	crt, key, err := client.GetCert()
	if err != nil {
		return err
	}

	///////////////
	change := false
	//read old .crt
	if string(old_crt_content) != crt {
		change = true
	}

	//read old .key
	if string(old_key_content) != key {
		change = true
	}

	//update the file
	if change {
		crt_file_err := file.FileOverwrite(c.Crt_path, crt)
		if crt_file_err != nil {
			return crt_file_err
		}

		key_file_err := file.FileOverwrite(c.Key_path, key)
		if key_file_err != nil {
			return key_file_err
		}

		if success_callback != nil {
			success_callback(crt, key)
		}
	}

	return nil
}
