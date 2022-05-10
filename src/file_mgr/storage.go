package file_mgr

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"github.com/meson-network/peer-node/src/storage_mgr"
)

//url_hash to a rel_path to "public" folder
func UrlHashToPublicFileRelPath(url_hash string) string {
	return filepath.Join(url_hash[0:4], url_hash[4:8], url_hash[8:12], url_hash[12:32])
}

func RequestPublicFile(url_hash string) (string, error) {

	file, err := GetFile(url_hash, true, true)
	if err != nil {
		return "", err
	}

	if file == nil {
		return "", errors.New("no such file")
	}

	SetLastReqTime(url_hash)

	if file.Status != STATUS_DOWNLOADED {
		return "", errors.New("file not downloaded")
	}

	abs_file_path, abs_file_path_err := storage_mgr.GetInstance().FileExist("file", "public", file.Rel_path)
	if abs_file_path_err != nil {
		return "", errors.New("file not exist on disk")
	}

	return abs_file_path, nil
}

func SetLastReqTime(url_hash string) {

	key := "last_req_time" + url_hash

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
		UpdateFile(map[string]interface{}{"last_req_unixtime": unixtime_now}, url_hash)
	}

	reference_plugin.GetInstance().Set(key, &unixtime_now, 1800)
}

func GetLastReqTime(url_hash string) (int64, error) {
	key := "last_req_time" + url_hash
	ref, _ := reference_plugin.GetInstance().Get(key)
	if ref != nil {
		return *(ref.(*int64)), nil
	}

	//get from sqlite
	file, file_err := GetFile(url_hash, true, true)
	if file_err != nil {
		return 0, file_err
	}

	if file == nil {
		return 0, errors.New("no such file")
	}

	reference_plugin.GetInstance().Set(key, file.Last_req_unixtime, 1800)
	return file.Last_req_unixtime, nil

}
