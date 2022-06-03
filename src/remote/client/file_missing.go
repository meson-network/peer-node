package client

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/tools/http/api"
	commonApi "github.com/meson-network/peer_common/api"
	"github.com/meson-network/peer_common/cached_file"
)

func FileMissing(fileHash string) {
	postData := &cached_file.Msg_Req_FileMissing{
		Missing_files: []string{fileHash},
	}
	result := &commonApi.API_META_STATUS{}
	err := api.POST_(EndPoint+"/api/node/file/missing", Token, postData, 30, result)
	if err != nil {
		basic.Logger.Errorln("FileMissing post err:", err, "fileHash:", fileHash)
	}

	if result.Meta_status <= 0 {
		basic.Logger.Errorln("FileMissing post err:", result.Meta_message, "fileHash:", fileHash)
	}
}
