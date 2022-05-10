package file_mgr

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/coreservice-io/utils/hash_util"
	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"github.com/meson-network/peer-node/src/storage_mgr"
)

func UrlToPublicFileHash(url string) string {
	return hash_util.SHA256String(url)[0:32]
}

func UrlToPublicFileRelPath(url string) string {
	nameHash := UrlToPublicFileHash(url)
	return filepath.Join(nameHash[0:4], nameHash[4:8], nameHash[8:12], nameHash[12:32])
}

func RequestPublicFile(file_hash string) (string, error) {

	file, err := GetFile(file_hash, true, true)
	if err != nil {
		return "", err
	}

	if file == nil {
		return "", errors.New("no such file")
	}

	SetLastReqTime(file_hash)

	if file.Status != STATUS_DOWNLOADED {
		return "", errors.New("file not downloaded")
	}

	abs_file_path, abs_file_path_err := storage_mgr.GetInstance().FileExist("file", "public", file.Rel_path)
	if abs_file_path_err != nil {
		return "", errors.New("file not exist on disk")
	}

	return abs_file_path, nil
}

func SetLastReqTime(file_name_hash string) {

	key := "last_req_time" + file_name_hash

	unixtime_now := time.Now().Unix()
	ref, _ := reference_plugin.GetInstance().Get(key)
	set_todb := false
	if ref != nil {
		if *(ref.(*int64)) < unixtime_now-5 {
			set_todb = true
		}
	} else {
		set_todb = true
	}

	if set_todb {
		UpdateFile(map[string]interface{}{"last_req_unixtime": unixtime_now}, file_name_hash)
	}

	reference_plugin.GetInstance().Set(key, &unixtime_now, 1800)
}

func GetLastReqTime(file_name_hash string) (int64, error) {
	key := "last_req_time" + file_name_hash
	ref, _ := reference_plugin.GetInstance().Get(key)
	if ref != nil {
		return *(ref.(*int64)), nil
	}

	//get from sqlite
	file, file_err := GetFile(file_name_hash, true, true)
	if file_err != nil {
		return 0, file_err
	}

	if file == nil {
		return 0, errors.New("no such file")
	}

	reference_plugin.GetInstance().Set(key, file.Last_req_unixtime, 1800)
	return file.Last_req_unixtime, nil

}
