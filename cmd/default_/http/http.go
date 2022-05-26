package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd/default_/http/api"
	"github.com/meson-network/peer-node/plugin/echo_plugin"
	"github.com/meson-network/peer-node/src/file_mgr"
	pErr "github.com/meson-network/peer-node/tools/errors"
	"github.com/meson-network/peer_common"
)

//httpServer example
func StartDefaultHttpSever() {
	httpServer := echo_plugin.GetInstance()
	api.ConfigApi(httpServer)
	api.DeclareApi(httpServer)

	//for handling private storage
	//httpServer.GET("/_personal_/*", func(ctx echo.Context) error {
	//	//storage_mgr.GetInstance()
	//	return ctx.HTML(http.StatusOK, "personal data")
	//})

	//for handling public storage
	// https://spec00-xxsdfsdffsdf-06-pzname.xxx.com/path1/path2/path3/1.jpg
	httpServer.GET("/*", func(ctx echo.Context) error {
		//access_token := ctx.QueryParam("access_token")
		//if access_token == "" {
		//	return ctx.HTML(http.StatusOK, "request is forbidden")
		//}

		//todo validate access_token
		//get random key in path
		v := strings.Split(ctx.Request().RequestURI, peer_common.MesonAccessTokenMark)
		accessToken := ""
		if len(v) == 1 {
			cookie, err := ctx.Cookie(peer_common.MesonAccessTokenMark)
			if err != nil {
				ctx.Error(errors.New("invalid access key"))
				return nil
			}
			accessToken = cookie.Value
		} else if len(v) == 2 {
			accessToken = v[1]
		} else {
			//return c.String(http.StatusUnauthorized, "invalid random key")
			ctx.Error(errors.New("invalid access key"))
			return nil
		}
		_ = accessToken
		//check random key
		//if !randomKeyMgr.GetInstance().CheckRandomKey(randKey) {
		//	//return c.String(http.StatusUnauthorized, "invalid random key")
		//	c.Error(errors.New("invalid access key"))
		//	return nil
		//}

		//set cookie
		//cookie := new(http.Cookie)
		//cookie.Name = peer_common.MesonAccessTokenMark
		//cookie.Value, _ = randomKeyMgr.GetInstance().GetRandomKey()
		//ctx.SetCookie(cookie)

		//get fileName
		fileName := v[0][1:]
		//check fileName legal

		ctx.Set("fileName", fileName)

		//get bindName
		//bindName := parseBindName(c.Request().Host)
		//check bindName legal

		//ctx.Set("bindName", bindName)

		//get fileHash
		//fileHash := bindName + fileName //to hash
		//ctx.Set("fileHash", fileHash)

		file_hash := ctx.Request().Header.Get("file_hash")
		if file_hash == "" {
			return ctx.HTML(http.StatusOK, "file_hash not defined")
		}

		//basic.Logger.Infoln(ctx.Request().URL)
		//basic.Logger.Infoln(file_mgr.UrlToPublicFileHash(ctx.Request().RequestURI))
		//basic.Logger.Infoln(file_mgr.UrlToPublicFileRelPath(ctx.Request().RequestURI))

		file_abs, file_header_json, file_abs_err := file_mgr.RequestPublicFile(file_hash)
		if file_abs_err != nil {
			errCode, err := pErr.ResolveStatusError(file_abs_err)
			basic.Logger.Debugln("file_mgr.RequestPublicFile errCode:", errCode, "err:", err)

			switch errCode {
			case -10001:
			case -10002:
			case -10003:
			case -10004:
			case -10005:
			case -10006:
			case -10007:
			}

			//todo file missing
			//redirect to server

			return ctx.HTML(404, "file not found")
		}

		//basic.Logger.Infoln("file_abs", file_abs)
		//basic.Logger.Infoln("file_header_json", file_header_json)

		for k, v := range file_header_json {
			for _, item := range v {
				ctx.Response().Header().Add(k, item)
			}
		}

		//todo get ignoreHeader
		ignoreHeader := map[string]struct{}{}
		err := httpServer.FileWithPause(ctx, file_abs, file_header_json, ignoreHeader)

		return err
	})

	err := httpServer.Start()
	if err != nil {
		basic.Logger.Fatalln(err)
	}
}

func CheckDefaultHttpServerStarted() bool {
	return echo_plugin.GetInstance().CheckStarted()
}

func ServerReloadCert() error {
	return echo_plugin.GetInstance().ReloadCert()
}
