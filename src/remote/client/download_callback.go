package client

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/tools/http/api"
	commonApi "github.com/meson-network/peer_common/api"
	"github.com/meson-network/peer_common/download"
)

func SuccessCallback(filehash string, file_size int64) {
	postData := &download.Msg_Req_Download_Callback_Success{
		//Origin_url: "",
		File_hash: filehash,
		File_size: file_size,
	}
	result := &commonApi.API_META_STATUS{}
	err := api.POST_(EndPoint+"/api/node/download/success", Token, postData, 30, result)
	if err != nil {
		basic.Logger.Errorln("SuccessCallback post err:", err, "fileHash:", filehash)
		return
	}

	if result.Meta_status <= 0 {
		basic.Logger.Errorln("SuccessCallback post err:", result.Meta_message, "fileHash:", filehash)
		return
	}
}

func FailedCallback(filehash string, download_code int) {

}
