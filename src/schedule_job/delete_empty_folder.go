package schedule_job

import (
	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/src/file_mgr"
)

func DeleteEmptyFolder() {
	const jobName = "DeleteEmptyFolder"

	job.Start(
		//job process
		jobName,
		func() {
			file_mgr.LoopDeleteEmptyFolder()
		},
		//onPanic callback
		nil, //todo upload panic
		2,
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
