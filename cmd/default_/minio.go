package default_

import (
	"errors"
	"path/filepath"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/src/cert_mgr"
	minio "github.com/minio/minio/cmd"
)

func RunMinio() error {

	runStorage, err := configuration.Config.GetBool("storage", true)
	if err != nil {
		return errors.New("storage [bool] in config error," + err.Error())
	}
	if !runStorage {
		return nil
	}

	//read config
	//folder
	storage_folder, err := configuration.Config.GetString("storage_folder", "m_storage")
	if err != nil {
		return errors.New("storage_folder [string] in config error," + err.Error())
	}
	if storage_folder == "" {
		storage_folder = "m_storage"
	}
	absPath := ""
	if filepath.IsAbs(storage_folder) {
		absPath = storage_folder
	} else {
		absPath = path_util.ExE_Path(storage_folder)
	}

	crt := cert_mgr.GetInstance().Crt_path
	certFolder := filepath.Dir(crt)

	minio.Main([]string{"", "server", absPath, "--address", "localhost:8080", "--console-address", "localhost:8081", "--certs-dir", certFolder})
	return nil
}
