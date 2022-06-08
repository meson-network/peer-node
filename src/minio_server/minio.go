package minio_server

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/configuration"
	"github.com/meson-network/peer-node/src/cert_mgr"
	"github.com/meson-network/peer-node/src/remote/client"
	minio "github.com/minio/minio/cmd"
)

//type MinioServer struct {
//	StorageFolder string //folder abs path
//	CertFolder    string // cert file abs path
//	ApiPort       string
//	ConsolePort   string
//	Domain        string
//}
//
//var minio_server *MinioServer

var ApiPort string

func RunMinio() error {
	runStorage, err := configuration.Config.GetBool("storage", true)
	if err != nil {
		return errors.New("storage [bool] in config error," + err.Error())
	}
	if !runStorage {
		return nil
	}

	//read config
	//folder
	storage_folder, err := configuration.Config.GetString("storage_folder", "m_storage")
	if err != nil {
		return errors.New("storage_folder [string] in config error," + err.Error())
	}
	if storage_folder == "" {
		storage_folder = "m_storage"
	}
	storage_folder_abs_path := ""
	if filepath.IsAbs(storage_folder) {
		storage_folder_abs_path = storage_folder
	} else {
		storage_folder_abs_path = path_util.ExE_Path(storage_folder)
	}

	crt := cert_mgr.GetInstance().Crt_path
	certFolder := filepath.Dir(crt)

	apiPort, err := configuration.Config.GetInt("storage_api_port", 0)
	if err != nil {
		return errors.New("storage_api_port [int] in config error," + err.Error())
	}
	if apiPort <= 0 || apiPort > 65535 {
		return errors.New("api port error")
	}
	ApiPort = strconv.Itoa(apiPort)

	consolePort, err := configuration.Config.GetInt("storage_console_port", 0)
	if err != nil {
		return errors.New("storage_console_port [int] in config error," + err.Error())
	}
	if consolePort <= 0 || consolePort > 65535 {
		return errors.New("console port error")
	}

	//get domain from remote
	nodeDomain, err := client.GetNodeDomain()
	if err != nil {
		return errors.New("get node domain error," + err.Error())
	}

	password, err := configuration.Config.GetString("storage_password", "")
	if err != nil {
		return errors.New("storage_console_port [int] in config error," + err.Error())
	}
	if password == "" {
		return errors.New("storage password not exist")
	}

	os.Setenv("MINIO_ROOT_USER", "mesonadmin")
	os.Setenv("MINIO_ROOT_PASSWORD", password)

	basic.Logger.Infoln("storage path:", storage_folder_abs_path)
	basic.Logger.Infoln("--address:", nodeDomain+":"+strconv.Itoa(apiPort))
	basic.Logger.Infoln("--console-address:", ":"+strconv.Itoa(consolePort))

	minio.Main([]string{"", "server", storage_folder_abs_path, "--address", nodeDomain + ":" + strconv.Itoa(apiPort), "--console-address", ":" + strconv.Itoa(consolePort), "--certs-dir", certFolder})
	return nil
}

//func RunMinio() error {
//
//	runStorage, err := configuration.Config.GetBool("storage", true)
//	if err != nil {
//		return errors.New("storage [bool] in config error," + err.Error())
//	}
//	if !runStorage {
//		return nil
//	}
//
//	//read config
//	//folder
//	storage_folder, err := configuration.Config.GetString("storage_folder", "m_storage")
//	if err != nil {
//		return errors.New("storage_folder [string] in config error," + err.Error())
//	}
//	if storage_folder == "" {
//		storage_folder = "m_storage"
//	}
//	absPath := ""
//	if filepath.IsAbs(storage_folder) {
//		absPath = storage_folder
//	} else {
//		absPath = path_util.ExE_Path(storage_folder)
//	}
//
//	crt := cert_mgr.GetInstance().Crt_path
//	certFolder := filepath.Dir(crt)
//
//	//minio.Main([]string{"", "server", absPath, "--address", "localhost:8080", "--console-address", "localhost:8081", "--certs-dir", certFolder})
//	minio.Main([]string{"", "server", absPath, "--address", "spec00-fekcdikbhgkjhxx.mesontracking.com:20443", "--certs-dir", certFolder})
//	return nil
//}
