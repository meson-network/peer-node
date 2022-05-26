package cert

import (
	"errors"
	"io/ioutil"

	"github.com/coreservice-io/job"
	"github.com/coreservice-io/utils/hash_util"
	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/src/remote/client"
	error_tool "github.com/meson-network/peer-node/tools/errors"
	"github.com/meson-network/peer-node/tools/file"
	"github.com/meson-network/peer-node/tools/http/api"
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

	//todo use default path
	//todo if cert file not exist, create folder and cert file

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

//success_callback func(string crt, string key)
func (c *CertMgr) UpdateCert(success_callback func(string, string)) error {

	//check hash
	hashRes := &dns.Msg_Resp_CertHash{}
	err := api.Get(client.EndPoint+"/api/node/cert/hash", client.Token, hashRes)
	if err != nil {
		return err
	}
	if hashRes.Meta_status <= 0 {
		return errors.New(hashRes.Meta_message)
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
	if old_content_hash == hashRes.Hash {
		return nil
	}

	res := &dns.Msg_Resp_Cert{}
	err = api.Get(client.EndPoint+"/api/node/cert", client.Token, res)

	if err != nil {
		return err
	}

	if res.Meta_status <= 0 {
		return errors.New(res.Meta_message)
	}

	///////////////
	change := false
	//read old .crt
	if string(old_crt_content) != res.Crt {
		change = true
	}

	//read old .key
	if string(old_key_content) != res.Key {
		change = true
	}

	//update the file
	if change {
		crt_file_err := file.FileOverwrite(c.Crt_path, res.Crt)
		if crt_file_err != nil {
			return crt_file_err
		}

		key_file_err := file.FileOverwrite(c.Key_path, res.Key)
		if key_file_err != nil {
			return key_file_err
		}

		if success_callback != nil {
			success_callback(res.Crt, res.Key)
		}
	}

	return nil
}

func (c *CertMgr) ScheduleUpdateJob(success_callback func(string, string)) {
	job.Start(
		"update_cert_job",
		// job process
		func() {
			c.UpdateCert(success_callback)
		},
		// onPanic callback, run if panic happened
		error_tool.PanicHandler,
		// job interval in seconds
		3600, //
		// job type
		// job.TYPE_PANIC_REDO  auto restart if panic
		// job.TYPE_PANIC_RETURN  stop if panic
		job.TYPE_PANIC_REDO,
		// check continue callback, the job will stop running if return false
		// the job will keep running if this callback is nil
		func(job *job.Job) bool {
			return true
		},
		// onFinish callback
		func(inst *job.Job) {

		},
	)
}
