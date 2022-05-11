package file_mgr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"github.com/meson-network/peer-node/src/storage_mgr"
)

//url_hash to a rel_path to "public" folder
func UrlHashToPublicFileRelPath(url_hash string) string {
	return filepath.Join(url_hash[0:4], url_hash[4:8], url_hash[8:16], url_hash)
}

//return abs_file_path file_header_json error
func RequestPublicFile(url_hash string) (string, map[string][]string, error) {

	file, err := GetFile(url_hash, true, true)
	if err != nil {
		return "", nil, err
	}

	if file == nil {
		return "", nil, errors.New("no such file")
	}

	SetLastReqTime(url_hash)

	if file.Status != STATUS_DOWNLOADED {
		return "", nil, errors.New("file not downloaded")
	}

	abs_file_path, abs_file_path_err := storage_mgr.GetInstance().FileExist("file", "public", file.Rel_path)
	if abs_file_path_err != nil {
		return "", nil, errors.New("file not exist on disk")
	}

	abs_file_header_path, abs_file_header_path_err := storage_mgr.GetInstance().FileExist("file", "public", file.Rel_path+".header")
	if abs_file_header_path_err != nil {
		return "", nil, errors.New("file header not exist on disk")
	}

	hfile, hfile_err := ioutil.ReadFile(abs_file_header_path)
	if hfile_err != nil {
		return "", nil, errors.New("file header read error")
	}

	header_json := make(map[string][]string)
	header_json_err := json.Unmarshal([]byte(hfile), &header_json)
	if header_json_err != nil {
		return "", nil, errors.New("file header json error")
	}

	return abs_file_path, header_json, nil
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
