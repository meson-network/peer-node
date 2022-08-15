package download_mgr

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/coreservice-io/utils/path_util"
	"github.com/imroc/req"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/cdn_cache_folder"
	"github.com/meson-network/peer-node/src/file_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
	"github.com/meson-network/peer-node/tools/file"
)

//todo put these consts into common as will be used for both server and peer
const NODE_DOWNLOAD_CODE_ERR = -10001                   //general download failure
const NODE_DOWNLOAD_CODE_ERR_BUSY = -10002              //active shutdown cause of max_downloaders limited , system too busy
const NODE_DOWNLOAD_CODE_ERR_SLOW = -10003              //active shutdown cause of too slow download at high traffic time
const NODE_DOWNLOAD_CODE_ERR_OTHER_DOWNLOADING = -10004 //active shutdown cause of someone else is downloading it
const NODE_DOWNLOAD_CODE_ERR_OVERSIZE = -10005          //active shutdown cause of single file size limit
const NODE_DOWNLOAD_CODE_ERR_DISK_SPACE = -10006        //active shutdown cause of single file size limit

const max_downloaders = 10
const max_file_size_bytes = 1024 * 1024 * 1024 //1GB limit
const min_speed_byte_per_sec = 1024 * 250      //active shutdown if reach (max_downloaders*70%) and download speed is below 250kb/second sec
const safe_seconds = 20

var total_downloaders int64

func clean_download(filehash string, file_path string) {
	os.Remove(file_path)
	file_mgr.DeleteFile(filehash)
}

func PreCheckTask(remoteUrl string, size_limt int64) error {
	if GetTotalDownloaderNum() >= max_downloaders {
		return pErr.NewStatusError(NODE_DOWNLOAD_CODE_ERR_BUSY, "too many running download task")
	}

	//check space
	freeSize := cdn_cache_folder.GetInstance().GetFreeSize()
	if freeSize < cdn_cache_folder.FreeSpaceLine {
		return pErr.NewStatusError(NODE_DOWNLOAD_CODE_ERR_DISK_SPACE, "have not enough space")
	}

	//try to check size
	r := req.New()
	r.SetTimeout(time.Duration(15) * time.Second)
	result, err := r.Get(remoteUrl)
	if err != nil {
		return pErr.NewStatusError(NODE_DOWNLOAD_CODE_ERR, "request origin err")
	}
	if result == nil {
		return pErr.NewStatusError(NODE_DOWNLOAD_CODE_ERR, "request origin err")
	}
	defer func() {
		if result.Response().Body != nil {
			result.Response().Body.Close()
		}
	}()

	value, exist := result.Response().Header["Content-Length"]
	if exist && len(value) > 0 {
		size, err := strconv.Atoi(value[0])
		if err == nil && size > 0 {
			if int64(size) > size_limt {
				return pErr.NewStatusError(NODE_DOWNLOAD_CODE_ERR_OVERSIZE, "file too big")
			}
			if int64(size) > freeSize {
				return pErr.NewStatusError(NODE_DOWNLOAD_CODE_ERR_DISK_SPACE, "have not enough space")
			}
		}
	}

	return nil
}

func StartDownloader(
	remoteUrl string,
	file_hash string,
	no_access_maintain_sec int64,
	size_limit int64,
	callback_succeed func(filehash string, file_size int64),
	callback_failed func(filehash string, download_code int),
) {

	file_relpath := file_mgr.UrlHashToPublicFileRelPath(file_hash)
	des_path := path.Join(cdn_cache_folder.GetInstance().GetCacheFileSaveFolderPath(), file_relpath)

	old_file, file_err := file_mgr.GetFile(file_hash, false, true)
	if file_err != nil {
		callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR)
		return
	}

	if old_file != nil {
		if old_file.Status == file_mgr.STATUS_DOWNLOADED {
			//check file exist on disk
			absPath := file_mgr.GetFileAbsPath(file_hash)
			exist, err := path_util.AbsPathExist(absPath)
			//file not exist on disk
			if err != nil || !exist {
				file_mgr.DeleteFile(file_hash)
			} else {
				callback_succeed(file_hash, old_file.Size_byte)
				return
			}
		} else {
			callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR_OTHER_DOWNLOADING)
			return
		}
	}

	////////system limit checker//////////
	atomic.AddInt64(&total_downloaders, 1)
	defer func() {
		atomic.AddInt64(&total_downloaders, -1)
	}()

	if total_downloaders >= max_downloaders {
		callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR_BUSY)
		return
	}
	///////////////////////////////////////

	nowTime := time.Now().UTC().Unix()
	file_mgr.CreateFile(&file_mgr.FileModel{
		File_hash:              file_hash,
		Last_req_unixtime:      nowTime,
		No_access_maintain_sec: no_access_maintain_sec,
		Size_byte:              0,
		Rel_path:               file_relpath,
		Status:                 file_mgr.STATUS_DOWNLOADING,
	})

	//dont forget to delete old fild otherwise you may append content after old content
	os.Remove(des_path)

	req, req_err := grab.NewRequest(des_path, remoteUrl)
	if req_err != nil {
		clean_download(file_hash, des_path)
		callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR)
		return
	}

	client := grab.NewClient()
	basic.Logger.Debugln("download from :", remoteUrl)
	basic.Logger.Debugln("download to :", des_path)

	resp := client.Do(req)
	if err := resp.Err(); err != nil {
		clean_download(file_hash, des_path)
		callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR)
		return
	}

	start_time := time.Now()
	t := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-t.C:
			//check size limits
			if resp.BytesComplete() > size_limit {
				clean_download(file_hash, des_path)
				callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR_OVERSIZE)
				return
			}

			//check too slow downloading
			elapsed := time.Since(start_time)
			if elapsed.Seconds() > safe_seconds && total_downloaders > (max_downloaders*0.7) && resp.BytesComplete() < (min_speed_byte_per_sec*10) {
				clean_download(file_hash, des_path)
				callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR_SLOW)
				return
			}

			nowTime = time.Now().UTC().Unix()
			file_mgr.UpdateFile(map[string]interface{}{
				"last_req_unixtime": nowTime,
				"size_byte":         resp.BytesComplete(),
			}, file_hash)

		case <-resp.Done:
			if resp.Err() != nil {
				clean_download(file_hash, des_path)
				callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR)

			} else {

				//save header
				hj, hj_err := json.Marshal(resp.HTTPResponse.Header)
				if hj_err != nil {
					clean_download(file_hash, des_path)
					callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR)
					return
				}

				h_file_err := file.FileOverwrite(des_path+".header", string(hj))
				if h_file_err != nil {
					clean_download(file_hash, des_path)
					callback_failed(file_hash, NODE_DOWNLOAD_CODE_ERR)
					return
				}

				//update database
				nowTime = time.Now().UTC().Unix()
				file_mgr.UpdateFile(map[string]interface{}{
					"last_req_unixtime": nowTime,
					"size_byte":         resp.BytesComplete(),
					"status":            file_mgr.STATUS_DOWNLOADED,
				}, file_hash)
				callback_succeed(file_hash, resp.BytesComplete())
				cdn_cache_folder.GetInstance().AddCacheUsedSize(resp.BytesComplete())
			}
			return
		}
	}

}

func GetTotalDownloaderNum() int64 {
	return total_downloaders
}
