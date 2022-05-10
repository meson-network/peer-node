package download_mgr

import "github.com/meson-network/peer-node/src/remote/download"

func DoTask(callback_succeed func(filehash string, file_local_abs_path string),
	callback_failed func(filehash string, download_code int)) error {

	dt, dt_err := download.GetDownloadTask()
	if dt_err != nil {
		return dt_err
	}

	StartDownloader(dt.Origin_url, dt.Url_hash, callback_succeed, callback_failed)
	return nil
}
