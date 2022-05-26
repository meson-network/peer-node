package echo_plugin

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func (hs *EchoServer) SetPauseSeconds(secs int64) {
	hs.PauseMoment = time.Now().Unix() + secs
}

func (hs *EchoServer) GetPauseMoment() int64 {
	return hs.PauseMoment
}

func (hs *EchoServer) FileWithPause(c echo.Context, filePath string, header map[string][]string, ignoreHeaderMap map[string]struct{}) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return echo.NotFoundHandler(c)
	}
	defer f.Close()
	fi, _ := f.Stat()

	for headerKey, headerValue := range header {
		_, exist := ignoreHeaderMap[headerKey]
		if exist {
			continue
		}
		for _, v := range headerValue {
			c.Response().Header().Add(headerKey, v)
		}
	}
	ServeContent(hs, c.Response(), c.Request(), fi.Name(), fi.ModTime(), f)
	return
}
