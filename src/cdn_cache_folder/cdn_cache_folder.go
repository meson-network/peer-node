package cdn_cache_folder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer_common/cdn_cache"
)

const SpaceHolderFolder = "spaceholder"
const CacheFileFolder = "files"

const eachHoldFileSize = 200 * 1024 * 1024

type CdnCacheFolder struct {
	Abs_path           string //absolute folder path of the root folder
	Cache_provide_size int64  //size cdn cache can use in byte
	Cache_used_size    int64  //size cdn cache used in byte

	sizeLock sync.RWMutex
}

var cdn_cache_mgr_pointer *CdnCacheFolder

func Init() error {
	if cdn_cache_mgr_pointer != nil {
		return nil
	}

	toml_conf := conf.Get_config().Toml_config

	if toml_conf.Cache.Size <= 0 {
		return errors.New("[cache.size_GB] not configured correctly")
	}

	if toml_conf.Cache.Size < cdn_cache.MIN_CACHE_SIZE {
		return fmt.Errorf("[cache.size_GB] must be at least %d GB", cdn_cache.MIN_CACHE_SIZE)
	}

	//cdn_cache dir
	cacheFolder := toml_conf.Cache.Folder
	if cacheFolder == "" { //set a default value
		cacheFolder = "m_cache"
		//return errors.New("[cache.folder] not configured correctly")
	}

	absPath := ""
	if filepath.IsAbs(cacheFolder) {
		absPath = cacheFolder
	} else {
		absPath = path_util.ExE_Path(cacheFolder)
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
		return errors.New("cache_folder path is not a directory")
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
		Cache_provide_size: int64(toml_conf.Cache.Size) * 1024 * 1024 * 1024,
		Cache_used_size:    0,
	}

	return nil
}

func GetInstance() *CdnCacheFolder {
	return cdn_cache_mgr_pointer
}

func (cf *CdnCacheFolder) GetCacheFileSaveFolderPath() string {
	return filepath.Join(cf.Abs_path, CacheFileFolder)
}

func (cf *CdnCacheFolder) SyncCacheFolderSize() {
	//var size int64
	var size struct {
		TotalSize int64 `json:"total_size"`
	}
	err := sqlite_plugin.GetInstance().Table("file").Select("sum(size_byte) as total_size").Where("status='DOWNLOADED'").Take(&size).Error
	if err != nil {
		if !strings.Contains(err.Error(), "converting NULL to int64 is unsupported") {
			basic.Logger.Errorln("syncCacheFolderSize err:", err)
		}
		return
	}
	//basic.Logger.Infoln(size)

	cf.SetCacheUsedSize(size.TotalSize)
}
