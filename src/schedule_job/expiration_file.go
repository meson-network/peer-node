package schedule_job

import (
	"os"
	"time"

	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/file_mgr"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/cached_file"
)

const expireTime = 3600 * 6 //no access in 6 hours

func ScanExpirationFile() {
	const jobName = "ExpirationFile"

	job.Start(
		//job process
		jobName,
		func() {
			reportExpiredFiles()
		},
		//onPanic callback
		nil, //todo upload panic
		600,
		// job type
		// UJob.TYPE_PANIC_REDO  auto restart if panic
		// UJob.TYPE_PANIC_RETURN  stop if panic
		job.TYPE_PANIC_REDO,
		// check continue callback, the job will stop running if return false
		// the job will keep running if this callback is nil
		nil,
		// onFinish callback
		nil,
	)
}

func reportExpiredFiles() error {
	//todo get files no accessed
	nowTime := time.Now().UTC().Unix()
	timeLine := nowTime - expireTime
	offset := 0
	for {
		result, err := file_mgr.QueryFile(nil, &timeLine, &[]string{file_mgr.STATUS_DOWNLOADED}, nil, 500, offset, false, false)
		if err != nil {
			return err
		}
		if len(result.Files) == 0 {
			return nil
		}

		offset += len(result.Files)

		expiredFiles := []string{}
		for _, v := range result.Files {
			expiredFiles = append(expiredFiles, v.File_hash)
		}

		//send to server
		postData := &cached_file.Msg_Req_FileExpire{
			Expired_files: expiredFiles,
		}
		res := &cached_file.Msg_Resp_FileExpire{}
		err = api.POST_(client.EndPoint+"/api/node/file/expire", client.Token, postData, 30, res)
		if err != nil {
			basic.Logger.Errorln("reportExpiredFiles post error:", err)
			continue
		}

		if res.Meta_status <= 0 {
			basic.Logger.Errorln("reportExpiredFiles post error:", res.Meta_message)
			continue
		}

		for _, v := range result.Files {
			//todo delete file and header on disk
			absPath := file_mgr.GetFileAbsPath(v.File_hash)
			os.Remove(absPath)
			os.Remove(absPath + ".header")
			file_mgr.DeleteFile(v.File_hash)
		}
	}
}
