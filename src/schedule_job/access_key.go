package schedule_job

import (
	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/src/access_key_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
)

func RenewAccessKey() {
	const jobName = "RenewAccessKey"

	job.Start(
		//job process
		jobName,
		func() {
			access_key_mgr.GetInstance().GenNewRandomKey()
		},
		//onPanic callback
		pErr.PanicHandler, //todo upload panic
		300,
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
