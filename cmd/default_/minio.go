package default_

import (
	"errors"
	"path/filepath"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/configuration"
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
	storage_folder, err := configuration.Config.GetString("storage_folder", "./storage")
	if err != nil {
		return errors.New("storage_folder [string] in config error," + err.Error())
	}
	if storage_folder == "" {
		storage_folder = "./storage"
	}
	absPath := ""
	if filepath.IsAbs(storage_folder) {
		absPath = storage_folder
	} else {
		absPath = path_util.ExE_Path(storage_folder)
	}

	minio.Main([]string{"", "server", absPath, "--address", ":8080", "--console-address", ":8081"})
	return nil
}
