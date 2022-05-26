package schedule_job

import (
	"github.com/coreservice-io/job"
)

func ScheduleCheckVersion() {
	const jobName = "CheckVersion"

	job.Start(
		//job process
		jobName,
		func() {
			//heartbeat.SendHeartBeat()
		},
		//onPanic callback
		nil, //todo upload panic
		5,
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
