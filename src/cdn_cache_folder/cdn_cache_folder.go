package cdn_cache_folder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer_common/cdn_cache"
)

const SpaceHolderFolder = "spaceholder"
const CacheFileFolder = "files"

const eachHoldFileSize = 200 * 1024 * 1024
const FreeSpaceLine = 1 * 1024 * 1024 * 1024 //1GB

const defaultCdnCacheSizeGB = 40

type CdnCacheFolder struct {
	Abs_path string //absolute folder path of the root folder
	//Total_size_config  int    //total size in G bytes used for whole cdn_cache space
	Cache_provide_size int64 //size cdn cache can use in byte
	Cache_used_size    int64 //size cdn cache used in byte

	sizeLock sync.RWMutex

	//Private_size   int    //size in G bytes used for node's owner personal uploader space
}

var cdn_cache_mgr_pointer *CdnCacheFolder

func Init() error {
	if cdn_cache_mgr_pointer != nil {
		return nil
	}

	cdn_cache_size, cdn_cache_size_err := configuration.Config.GetInt("cdn_cache_size", defaultCdnCacheSizeGB)
	if cdn_cache_size_err != nil || cdn_cache_size == 0 {
		return errors.New("cdn_cache_size not configured correctly")
	}

	if cdn_cache_size < cdn_cache.MIN_CACHE_SIZE {
		return fmt.Errorf("cdn_cache_size must be at least %d GB", cdn_cache.MIN_CACHE_SIZE)
	}

	//cdn_cache dir
	sf, sf_err := configuration.Config.GetString("cdn_cache_folder", "m_cache")
	if sf_err != nil {
		return errors.New("cdn_cache_folder not configured correctly")
	}

	absPath := ""
	if filepath.IsAbs(sf) {
		absPath = sf
	} else {
		absPath = path_util.ExE_Path(sf)
	}

	err := os.MkdirAll(absPath, 0777)
	if err != nil {
		return fmt.Errorf("create meson_cdn_cache folder error:%s", err)
	}
	err = os.Chmod(absPath, 0777)
	if err != nil {
		return fmt.Errorf("modify meson_cdn_cache folder permission error:%s", err)
	}

	//check path is a dir
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("cdn_cache_folder path is not a directory")
	}

	//err = dm.checkAllFileInDb(onFileMissing)
	//if err != nil {
	//	return err
	//}
	//
	//dm.continueDownloadingFiles()
	//
	////bgJob
	//dm.loopScanLeakFiles()
	//dm.loopDeleteEmptyFolder()
	//dm.loopRetryDownloadNotify()

	cdn_cache_mgr_pointer = &CdnCacheFolder{
		Abs_path:           absPath,
		Cache_provide_size: int64(cdn_cache_size) * 1024 * 1024 * 1024,
		Cache_used_size:    0,
		//Private_size:   int(float64(storage_size) * cdn_cache.MAX_STOR_PERSONAL_RATIO),
	}

	return nil
}

func GetInstance() *CdnCacheFolder {
	return cdn_cache_mgr_pointer
}

func (cf *CdnCacheFolder) GetCacheFileSaveFolderPath() string {
	return filepath.Join(cf.Abs_path, CacheFileFolder)
}
