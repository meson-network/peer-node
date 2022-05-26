package file_mgr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/plugin/reference_plugin"
	pErr "github.com/meson-network/peer-node/tools/errors"
)

//file_hash to a rel_path to "public" folder
func UrlHashToPublicFileRelPath(file_hash string) string {
	return filepath.Join(file_hash[0:4], file_hash[4:8], file_hash[8:16], file_hash)
}

//return abs_file_path file_header_json error
func RequestPublicFile(file_hash string) (string, map[string][]string, error) {

	file, err := GetFile(file_hash, true, true)
	if err != nil {
		return "", nil, pErr.NewStatusError(-10001, err.Error())
	}

	if file == nil {
		return "", nil, pErr.NewStatusError(-10002, "no such file")
	}

	if file.Status != STATUS_DOWNLOADED {
		return "", nil, pErr.NewStatusError(-10003, "file not downloaded")
	}

	absPath := GetFileAbsPath(file_hash)
	exist, err := path_util.AbsPathExist(absPath)
	if err != nil || !exist {
		return "", nil, pErr.NewStatusError(-10004, "file not exist on disk")
	}

	headerAbsPath := absPath + ".header"
	exist, err = path_util.AbsPathExist(headerAbsPath)
	if err != nil || !exist {
		return "", nil, pErr.NewStatusError(-10005, "file header not exist on disk")
	}

	hfile, hfile_err := ioutil.ReadFile(headerAbsPath)
	if hfile_err != nil {
		return "", nil, pErr.NewStatusError(-10006, "file header read error")
	}

	header_json := make(map[string][]string)
	header_json_err := json.Unmarshal([]byte(hfile), &header_json)
	if header_json_err != nil {
		return "", nil, pErr.NewStatusError(-10007, "file header json error")
	}

	SetLastReqTime(file_hash)

	return absPath, header_json, nil
}

func SetLastReqTime(file_hash string) {

	key := "last_req_time" + file_hash

	unixtime_now := time.Now().Unix()
	ref, _ := reference_plugin.GetInstance().Get(key)
	set_todb := false
	if ref != nil {
		if *(ref.(*int64)) < unixtime_now-5 { //todo 5s=>30s??
			set_todb = true
		}
	} else {
		set_todb = true
	}

	if set_todb {
		UpdateFile(map[string]interface{}{"last_req_unixtime": unixtime_now}, file_hash)
	}

	reference_plugin.GetInstance().Set(key, &unixtime_now, 1800)
}

func GetLastReqTime(file_hash string) (int64, error) {
	key := "last_req_time" + file_hash
	ref, _ := reference_plugin.GetInstance().Get(key)
	if ref != nil {
		return *(ref.(*int64)), nil
	}

	//get from sqlite
	file, file_err := GetFile(file_hash, true, true)
	if file_err != nil {
		return 0, file_err
	}

	if file == nil {
		return 0, errors.New("no such file")
	}

	reference_plugin.GetInstance().Set(key, file.Last_req_unixtime, 1800)
	return file.Last_req_unixtime, nil

}
