package schedule_job

import (
	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/src/cert_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
	minio "github.com/minio/minio/cmd"
)

func UpdateCert() {
	const jobName = "UpdateCert"

	job.Start(
		//job process
		jobName,
		func() {
			err := cert_mgr.GetInstance().UpdateCert(func(crt, key string) {
				err := http.ServerReloadCert()
				if err != nil {
					basic.Logger.Errorln("schedule UpdateCert http.ServerReloadCert error:", err)
				}
				basic.Logger.Debugln("minio.GetCertManger().ReloadCerts()")
				im := minio.GetCertManger()
				if im != nil {
					im.ReloadCerts()
				}
			})
			if err != nil {
				basic.Logger.Errorln("schedule UpdateCert error:", err)
			}
		},
		//onPanic callback
		pErr.PanicHandler, //todo upload panic
		3600,              //todo 3600 in production
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
