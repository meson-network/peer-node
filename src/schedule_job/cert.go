package schedule_job

import (
	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/src/cert_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
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
			})
			if err != nil {
				basic.Logger.Errorln("schedule UpdateCert error:", err)
			}
		},
		//onPanic callback
		pErr.PanicHandler, //
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
