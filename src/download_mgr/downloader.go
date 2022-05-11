package download_mgr

import (
	"encoding/json"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/file_mgr"
	"github.com/meson-network/peer-node/src/storage_mgr"
	"github.com/meson-network/peer-node/tools/file"
)

//todo put these consts into common as will be used for both server and peer
const NODE_DOWNLOAD_CODE_ERR = -1                   //general download failure
const NODE_DOWNLOAD_CODE_ERR_BUSY = -2              //active shutdown cause of max_downloaders limited , system too busy
const NODE_DOWNLOAD_CODE_ERR_SLOW = -3              //active shutdown cause of too slow download at high traffic time
const NODE_DOWNLOAD_CODE_ERR_OTHER_DOWNLOADING = -4 //active shutdown cause of someone else is downloading it
const NODE_DOWNLOAD_CODE_ERR_OVERSIZE = -5          //active shutdown cause of single file size limit

const max_downloaders = 10
const max_file_size_bytes = 1024 * 1024 * 1024 //1GB limit
const min_speed_byte_per_sec = 1024 * 250      //active shutdown if reach (max_downloaders*70%) and download speed is below 250kb/second sec

var total_downloaders int64

func clean_download(filehash string, file_path string) {
	os.Remove(file_path)
	file_mgr.DeleteFile(filehash)
}

func StartDownloader(
	remoteUrl string,
	url_hash string,
	callback_succeed func(filehash string, file_local_abs_path string),
	callback_failed func(filehash string, download_code int),
) {

	file_relpath := file_mgr.UrlHashToPublicFileRelPath(url_hash)

	des_path := path.Join(storage_mgr.GetInstance().Storage_folder, "file", "public", file_relpath)

	old_file, file_err := file_mgr.GetFile(url_hash, true, true)
	if file_err != nil {
		callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR)
		return
	}

	if old_file != nil {
		if old_file.Status == file_mgr.STATUS_DOWNLOADED {
			callback_succeed(url_hash, des_path)
		} else {
			callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR_OTHER_DOWNLOADING)
		}
		return
	}

	////////system limit checker//////////
	atomic.AddInt64(&total_downloaders, 1)
	defer func() {
		atomic.AddInt64(&total_downloaders, -1)
	}()

	if total_downloaders >= max_downloaders {
		callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR_BUSY)
		return
	}
	///////////////////////////////////////

	client := grab.NewClient()
	basic.Logger.Debugln("download from :", remoteUrl)
	basic.Logger.Debugln("download to :", des_path)

	file_mgr.CreateFile(&file_mgr.FileModel{
		Url_hash:               url_hash,
		Last_req_unixtime:      time.Now().Unix(),
		Last_scan_unixtime:     time.Now().Unix(),
		Last_download_unixtime: time.Now().Unix(),
		Size_byte:              0,
		Rel_path:               file_relpath,
		Status:                 file_mgr.STATUS_DOWNLOADING,
		Type:                   file_mgr.TYPE_PUBLIC,
	})

	req, req_err := grab.NewRequest(des_path, remoteUrl)
	if req_err != nil {
		clean_download(url_hash, des_path)
		callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR)
		return
	}

	start_time := time.Now()
	t := time.NewTicker(2 * time.Second)

	resp := client.Do(req)

	for {
		select {
		case <-t.C:
			//check size limits
			if resp.BytesComplete() > max_file_size_bytes {
				clean_download(url_hash, des_path)
				callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR_OVERSIZE)
				return
			}

			//check too slow downloading
			elapsed := time.Since(start_time)
			if elapsed.Seconds() > 10 && total_downloaders > (max_downloaders*0.7) && resp.BytesComplete() < (min_speed_byte_per_sec*10) {
				clean_download(url_hash, des_path)
				callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR_SLOW)
				return
			}

			file_mgr.UpdateFile(map[string]interface{}{
				"last_req_unixtime":      time.Now().Unix(),
				"last_scan_unixtime":     time.Now().Unix(),
				"last_download_unixtime": time.Now().Unix(),
				"size_byte":              resp.BytesComplete(),
			}, url_hash)

		case <-resp.Done:
			if resp.Err() != nil {
				clean_download(url_hash, des_path)
				callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR)

			} else {

				//save header
				hj, hj_err := json.Marshal(resp.HTTPResponse.Header)
				if hj_err != nil {
					//basic.Logger.Errorln(hj_err)
					clean_download(url_hash, des_path)
					callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR)
					return
				}

				h_file_err := file.FileOverwrite(des_path+".header", string(hj))
				if h_file_err != nil {
					//basic.Logger.Errorln(h_file_err)
					clean_download(url_hash, des_path)
					callback_failed(url_hash, NODE_DOWNLOAD_CODE_ERR)
					return
				}

				//update database
				file_mgr.UpdateFile(map[string]interface{}{
					"last_req_unixtime":      time.Now().Unix(),
					"last_scan_unixtime":     time.Now().Unix(),
					"last_download_unixtime": time.Now().Unix(),
					"size_byte":              resp.BytesComplete(),
					"status":                 file_mgr.STATUS_DOWNLOADED,
				}, url_hash)
				callback_succeed(url_hash, des_path)
			}
			return
		}
	}

}

func GetTotalDownloaderNum() int64 {
	return total_downloaders
}
