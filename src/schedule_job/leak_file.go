package schedule_job

import (
	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/src/file_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
)

func ScanLeakFile() {
	const jobName = "ScanLeakFile"

	job.Start(
		//job process
		jobName,
		func() {
			file_mgr.ScanLeakFiles()
		},
		//onPanic callback
		pErr.PanicHandler, //todo upload panic
		24*3600,
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
