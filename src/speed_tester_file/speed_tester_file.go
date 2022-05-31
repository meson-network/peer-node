package speed_tester_file

import (
	"math/rand"
	"os"
	"time"

	"github.com/coreservice-io/utils/path_util"
)

var fileName = "s_tester.bin"

func GetSpeedTesterFilePath() string {
	return path_util.ExE_Path(fileName)
}

func CheckTesterFile() error {
	absPath := path_util.ExE_Path(fileName)
	f, err := os.Stat(absPath)
	if err != nil || f.Size() != 32*1024*3200 { //100M
		genErr := genSpeedTesterFile()
		if genErr != nil {
			return genErr
		}
	}

	return nil
}

func genSpeedTesterFile() error {
	rand.Seed(time.Now().UnixNano())
	fileAbsPath := path_util.ExE_Path(fileName)
	var data = make([]byte, 32*1024, 32*1024) // Initialize an empty byte slice
	for i := 0; i < 32*1024; i++ {
		data[i] = byte(rand.Intn(255))
	}

	f, err := os.OpenFile(fileAbsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// 1600 round gen a 100M file
	for j := 0; j < 3200; j++ {
		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}
