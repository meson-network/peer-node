package file_request

import (
	"errors"
	"strings"
	"time"

	"github.com/coreservice-io/dns-common/spec00"
	"github.com/coreservice-io/safe_go"
	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/file_mgr"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/src/remote/file_missing"
	pErr "github.com/meson-network/peer-node/tools/errors"
	"github.com/meson-network/peer_common"
)

func HandleFileRequest(httpServer *echo_plugin.EchoServer) {
	// https://spec00-xxsdfsdffsdf-06-pzname.xxx.com/path1/path2/path3/1.jpg
	httpServer.GET("/*", func(ctx echo.Context) error {
		//todo validate access_token
		//get random key in path
		v := strings.Split(ctx.Request().RequestURI, peer_common.MesonAccessTokenMark)
		accessKey := ""
		if len(v) == 2 && v[1] != "" {
			accessKey = v[1]
		} else {
			//return c.String(http.StatusUnauthorized, "invalid random key")
			ctx.Error(errors.New("invalid access key"))
			return nil
		}
		//check random key
		clientToken, timeStamp, err := ParserAccessToken(accessKey)
		if err != nil {
			ctx.Error(errors.New("invalid access key"))
			return nil
		}
		if clientToken != client.Token {
			ctx.Error(errors.New("invalid access key"))
			return nil
		}
		if time.Now().UTC().Unix() > timeStamp+300 || time.Now().UTC().Unix() < timeStamp-300 { //+- 5 minutes
			ctx.Error(errors.New("invalid access key"))
			return nil
		}

		//get fileName
		fileName := v[0]
		//check fileName legal

		//ctx.Set("fileName", fileName)

		//get pullzone
		_, optionalStr, err := spec00.Parser(ctx.Request().Host)
		if err != nil {
			ctx.Error(errors.New("invalid request"))
			return nil
		}
		if len(optionalStr) == 0 {
			ctx.Error(errors.New("invalid request"))
			return nil
		}
		pullZone := optionalStr[0]

		//get fileHash
		fileHash := peer_common.GenFileHash(pullZone, fileName) //to hash
		//ctx.Set("fileHash", fileHash)

		basic.Logger.Debugln(ctx.Request().URL)
		basic.Logger.Debugln("pullzone:", pullZone)
		basic.Logger.Debugln("file:", fileName)
		basic.Logger.Debugln("hash:", fileHash)

		file_abs, file_header_json, file_abs_err := file_mgr.RequestPublicFile(fileHash)
		if file_abs_err != nil {
			errCode, err := pErr.ResolveStatusError(file_abs_err)
			basic.Logger.Errorln("file_mgr.RequestPublicFile errCode:", errCode, "err:", err)

			fileIsMissing := false
			switch errCode {
			case -10001: //get file info from db error
				//do nothing
			case -10002: //file not exist in db
				fileIsMissing = true
			case -10003: //file is downloading
				//do nothing
			case -10004: //file not exist on disk
				fileIsMissing = true
				//delete from db
				file_mgr.DeleteFile(fileHash)
			case -10005 - 10006 - 10007:
				//-10005 file header not exist
				//-10006 read file header error
				//-10007 unmarshal header error
				fileIsMissing = true
				//delete from db and disk
				file_mgr.DeleteFile(fileHash)
				file_mgr.RemoveFileFromDisk(fileHash)
			}

			//notify server file missing
			if fileIsMissing {
				safe_go.Go(func(args ...interface{}) {
					file_missing.FileMissing(fileHash)
				}, nil)
			}

			//redirect to server
			//https://pz-pullzone.meson.network/fileName
			redirectUrl := "https://pz-" + pullZone + ".meson.network" + fileName
			return ctx.Redirect(302, redirectUrl)
		}

		basic.Logger.Debugln("file_abs", file_abs)
		basic.Logger.Debugln("file_header_json", file_header_json)

		//todo get ignoreHeader
		ignoreHeader := map[string]struct{}{}
		err = httpServer.FileWithPause(ctx, file_abs, file_header_json, ignoreHeader)
		return err
	})
}

func ParserAccessToken(accessToken string) (nodeToken string, timeStamp int64, err error) {

	return "", 0, err
}
