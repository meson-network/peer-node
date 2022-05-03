package main

import (
	"os"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/cmd"
)

func main() {
	basic.InitLogger()

	//config app to run
	errRun := cmd.ConfigCmd().Run(os.Args)
	if errRun != nil {
		basic.Logger.Panicln(errRun)
	}
}
