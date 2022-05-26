package download

//import (
//	"errors"
//
//	"github.com/meson-network/peer-node/basic"
//	"github.com/meson-network/peer-node/src/remote/client"
//	"github.com/meson-network/peer-node/tools/http/api"
//	"github.com/meson-network/peer_common/download"
//)
//
//func GetDownloadTask() (*download.Msg_Resp_Download_Task, error) {
//
//	cient_, c_err := client.GetClient()
//	if c_err != nil {
//		basic.Logger.Fatalln(c_err)
//	}
//
//	res := &download.Msg_Resp_Download_Task{}
//	err := api.Get(cient_.EndPoint+"/api/node/download/task", cient_.Token, res)
//	if err != nil {
//		return nil, err
//	}
//
//	if res.Meta_status <= 0 {
//		return nil, errors.New(res.Meta_message)
//	}
//
//	return res, nil
//}
