package schedule_job

import (
	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/cmd/default_/http"
	"github.com/meson-network/peer-node/src/cert_mgr"
)

func UpdateCert() {
	const jobName = "UpdateCert"

	job.Start(
		//job process
		jobName,
		func() {
			cert_mgr.GetInstance().UpdateCert(func(crt, key string) {
				http.ServerReloadCert()
			})
		},
		//onPanic callback
		nil, //todo upload panic
		60,  //todo 3600 in production
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
