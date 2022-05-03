package log

import (
	"github.com/coreservice-io/log"
	"github.com/meson-network/peer-node/basic"
	"github.com/urfave/cli/v2"
)

func StartLog(clictx *cli.Context) {
	num := clictx.Int64("num")
	if num == 0 {
		num = 20
	}

	onlyerr := clictx.Bool("only_err")
	if onlyerr {
		basic.Logger.PrintLastN(num, []log.LogLevel{log.PanicLevel, log.FatalLevel, log.ErrorLevel})
	} else {
		basic.Logger.PrintLastN(num, []log.LogLevel{log.PanicLevel, log.FatalLevel, log.ErrorLevel, log.InfoLevel, log.WarnLevel, log.DebugLevel, log.TraceLevel})
	}
}
