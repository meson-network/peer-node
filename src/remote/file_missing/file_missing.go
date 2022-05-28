package file_missing

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/http/api"
	"github.com/meson-network/peer_common/cached_file"
)

func FileMissing(fileHash string) {
	postData := &cached_file.Msg_Req_FileMissing{
		Missing_files: []string{fileHash},
	}
	result := &cached_file.Msg_Resp_FileMissing{}
	err := api.POST_(client.EndPoint+"/api/node/file/missing", client.Token, postData, 30, result)
	if err != nil {
		basic.Logger.Errorln("FileMissing post err:", err, "fileHash:", fileHash)
	}

	if result.Meta_status <= 0 {
		basic.Logger.Errorln("FileMissing post err:", result.Meta_message, "fileHash:", fileHash)
	}
}
