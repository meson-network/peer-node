package schedule_job

import (
	"strconv"

	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/access_key_mgr"
	"github.com/meson-network/peer-node/src/node_info"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/src/version_mgr"
	"github.com/meson-network/peer_common/heart_beat"
)

func HeartBeat() {
	const jobName = "HeartBeat"

	job.Start(
		//job process
		jobName,
		func() {
			sendHeartBeat()
		},
		//onPanic callback
		nil, //todo upload panic
		30,
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

func sendHeartBeat() {
	accessKey, _ := access_key_mgr.GetInstance().GetRandomKey()
	portStr := echo_plugin.GetInstance().Http_port
	postData := &heart_beat.Msg_Req_HeartBeat{
		Node_id:    node_info.GetNodeId(),
		Port:       strconv.Itoa(portStr),
		Version:    version_mgr.NodeVersion,
		Access_key: accessKey,
	}
	client.SendHeartBeat(postData)
}
