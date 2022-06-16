package minio_server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/src/cert_mgr"
	"github.com/meson-network/peer-node/src/remote/client"
	minio "github.com/minio/minio/cmd"
)

var ApiPort string

func RunMinio() error {

	toml_conf := conf.Get_config().Toml_config
	if !toml_conf.Storage.Enable {
		return nil
	}

	basic.Logger.Infoln("Meson Storage enable, init node storage...")
	//read config
	//folder
	storage_folder := toml_conf.Storage.Folder
	if storage_folder == "" { //todo set a default value or return error
		storage_folder = "m_storage"
		//return errors.New("[storage.folder] not configured correctly")
	}
	storage_folder_abs_path := ""
	if filepath.IsAbs(storage_folder) {
		storage_folder_abs_path = storage_folder
	} else {
		storage_folder_abs_path = path_util.ExE_Path(storage_folder)
	}

	crt := cert_mgr.GetInstance().Crt_path
	certFolder := filepath.Dir(crt)

	apiPort := toml_conf.Storage.Api_port
	if apiPort <= 0 || apiPort > 65535 {
		return errors.New("api port error")
	}
	if apiPort == toml_conf.Https_port {
		return fmt.Errorf("storage api port [%d] already used in https port", apiPort)
	}
	ApiPort = strconv.Itoa(apiPort)

	consolePort := toml_conf.Storage.Console_port
	if consolePort <= 0 || consolePort > 65535 {
		return errors.New("console port error")
	}
	if consolePort == toml_conf.Https_port || consolePort == apiPort {
		return fmt.Errorf("storage api port [%d] already used in https port or api port", consolePort)
	}

	//get domain from remote
	nodeDomain, err := client.GetNodeDomain()
	if err != nil {
		return errors.New("get node domain error," + err.Error())
	}

	password := toml_conf.Storage.Password
	if password == "" {
		return errors.New("storage password not exist")
	}
	if len(password) < 6 {
		return errors.New("storage password length can not less than 6")
	}

	err = os.Setenv("MINIO_ROOT_USER", "mesonadmin")
	if err != nil {
		return errors.New("storage set env 'MINIO_ROOT_USER' error:" + err.Error())
	}
	err = os.Setenv("MINIO_ROOT_PASSWORD", password)
	if err != nil {
		return errors.New("storage set env 'MINIO_ROOT_PASSWORD' error:" + err.Error())
	}
	err = os.Setenv("MINIO_SERVER_URL", "https://"+nodeDomain+":"+strconv.Itoa(apiPort))
	if err != nil {
		return errors.New("storage set env 'MINIO_SERVER_URL' error:" + err.Error())
	}

	//basic.Logger.Infoln("storage path:", storage_folder_abs_path)
	//basic.Logger.Infoln("--address:", nodeDomain+":"+strconv.Itoa(apiPort))
	//basic.Logger.Infoln("--console-address:", ":"+strconv.Itoa(consolePort))

	basic.Logger.Infoln("Meson Storage api port:", apiPort)
	basic.Logger.Infoln("Meson Storage console port:", consolePort)
	basic.Logger.Infoln("Meson Storage console url:", "https://"+nodeDomain+":"+strconv.Itoa(apiPort))

	minio.Main([]string{"peer-node", "server", storage_folder_abs_path, "--address", ":" + strconv.Itoa(apiPort), "--console-address", ":" + strconv.Itoa(consolePort), "--certs-dir", certFolder})

	return nil
}
