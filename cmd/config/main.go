package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/meson-network/peer-node/configuration"
	"github.com/urfave/cli/v2"
)

func ConfigSetting(clictx *cli.Context) {

	configModify := false

	for _, v := range stringConfParams {
		if clictx.IsSet(v) {
			newValue := clictx.String(v)
			configuration.Config.Set(v, newValue)
			configModify = true
		}
	}

	for _, v := range float64ConfParams {
		if clictx.IsSet(v) {
			newValue := clictx.Float64(v)
			configuration.Config.Set(v, newValue)
			configModify = true
		}
	}

	for _, v := range boolConfPrams {
		if clictx.IsSet(v) {
			newValue := clictx.Bool(v)
			configuration.Config.Set(v, newValue)
			configModify = true
		}
	}

	//other custom flags example
	if clictx.IsSet("addpath") {
		folder, err := configuration.Config.GetProvideFolders()
		if err != nil {
			fmt.Println(err)
			return
		}

		folderPath := clictx.String("addpath")
		//check path legal
		//...

		for _, v := range folder {
			if v.AbsPath == folderPath {
				fmt.Println("Error: path already exist")
				return
			}
		}

		//input size
		var size int

		fmt.Printf("Please input provider folder size: ")
		_, err = fmt.Scanln(&size)
		if err != nil {
			fmt.Println("Read input size error: %s", err.Error())
			return
		}
		if size < 20 {
			fmt.Println("Error: minimum size is 20 GB")
			return
		}

		pf := configuration.ProvideFolder{
			AbsPath: folderPath,
			SizeGB:  size,
		}
		folder = append(folder, pf)
		configuration.SetProvideFolders(folder)
		configModify = true
		fmt.Println("new folder added:", folderPath, "size:", size, "GB")
	}

	if clictx.IsSet("removepath") {
		folder, err := configuration.Config.GetProvideFolders()
		if err != nil {
			fmt.Println(err)
			return
		}

		pathToRemove := clictx.String("removepath")

		removed := false
		for i, v := range folder {
			if v.AbsPath == pathToRemove {
				folder = append(folder[:i], folder[i+1:]...)
				removed = true
				break
			}
		}

		if removed {
			configuration.SetProvideFolders(folder)
			configModify = true
			fmt.Println("path removed:", pathToRemove)
		} else {
			fmt.Println("removepath failed, path not exist")
			return
		}
	}

	if configModify {
		err := configuration.Config.WriteConfig()
		if err != nil {
			color.Red("config save error:", err)
			return
		}
		fmt.Println("config modified")
		//fmt.Println(configuration.Config.GetConfigAsString())
	}
}
