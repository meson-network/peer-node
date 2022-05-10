package storage_mgr

import (
	"errors"
	"path/filepath"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer_common/storage"
)

type StorageMgr struct {
	Storage_folder string //absolute folder path of the root folder
	Total_size     int    //total size in G bytes used for whole storage space
	Private_size   int    //size in G bytes used for node's owner personal uploader space
}

var storage_mgr_pointer *StorageMgr

func Init() error {
	if storage_mgr_pointer != nil {
		return nil
	}

	sf, sf_err := configuration.Config.GetString("storage_folder", "")
	if sf_err != nil || sf == "" {
		return errors.New("storage_folder not configured correctly")
	}

	sf_absdir, abs_err := path_util.SmartExistPath(sf)
	if abs_err != nil {
		return errors.New(sf + " :storage_folder not exist , please reset your storage_folder ")
	}

	storage_size, storage_size_err := configuration.Config.GetInt("storage_size", 0)
	if storage_size_err != nil || storage_size == 0 {
		return errors.New("storage_size  not configured correctly")
	}

	if storage_size < storage.MIN_STOR_SIZE {
		return errors.New("storage_size must be at least 30 GB")
	}

	storage_mgr_pointer = &StorageMgr{
		Storage_folder: sf_absdir,
		Total_size:     storage_size,
		Private_size:   int(float64(storage_size) * storage.MAX_STOR_PERSONAL_RATIO),
	}

	return nil
}

func GetInstance() *StorageMgr {
	return storage_mgr_pointer
}

//check rel_path(folder/file) relative to storage folder exist
//return abs_path if exist
func (stor_m *StorageMgr) FileExist(rel_paths ...string) (string, error) {

	abs_path := filepath.Join(stor_m.Storage_folder, filepath.Join(rel_paths...))

	exsit, patherr := path_util.AbsPathExist(abs_path)
	if patherr != nil || !exsit {
		return "", errors.New("path not exist")
	}

	return abs_path, nil
}
