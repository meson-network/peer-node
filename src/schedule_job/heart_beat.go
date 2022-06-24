package schedule_job

import (
	"os"
	"strconv"

	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/basic/conf"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/access_key_mgr"
	"github.com/meson-network/peer-node/src/node_info"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/src/version_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
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
		pErr.PanicHandler, //todo upload panic
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

var isInitial = true

func sendHeartBeat() {
	accessKey, _ := access_key_mgr.GetInstance().GetRandomKey()
	portStr := echo_plugin.GetInstance().Http_port
	minio_apiPort := conf.Get_config().Toml_config.Storage.Api_port
	postData := &heart_beat.Msg_Req_HeartBeat{
		Node_id:      node_info.GetNodeId(),
		Port:         strconv.Itoa(portStr),
		Storage_port: strconv.Itoa(minio_apiPort),
		Version:      version_mgr.NodeVersion,
		Access_key:   accessKey,
		Initial:      isInitial,
	}

	result, err := client.SendHeartBeat(postData)
	if err != nil {
		basic.Logger.Errorln("SendHeartBeat err:", err)
		return
	}

	if isInitial {
		isInitial = false
	}
	switch result.Meta_status {
	case 1:
		//success
	case -10001: //post data error
		basic.Logger.Errorln("hb, post data error")
	case -10002: //heart request too fast
		basic.Logger.Errorln("hb, request too fast")
	case -10003: //version error
		basic.Logger.Errorln("hb, this version has expired, please download a new version")
		os.Exit(0)
	case -10004: //token error
		basic.Logger.Errorln("hb, token error, please set correct token in config")
	case -10005: //same ip exist
		basic.Logger.Errorln("hb, multiple nodes use the same ip")
		//os.Exit(0)
	case -10006: //ip can't resolve
		basic.Logger.Errorln("hb, ip resolve error")
	case -10007: //ip to spec00 host error
		basic.Logger.Errorln("hb, ip to host error")
	case -10008 - 10009: //ping back error
		basic.Logger.Errorln("hb, ping back error")
	case -10010 - 10011: //internal error
		basic.Logger.Errorln("hb, remote internal error")
	case -10099: //internal error
		basic.Logger.Errorln("hb, error code 10099 force stop")
		os.Exit(0)
	}
}
