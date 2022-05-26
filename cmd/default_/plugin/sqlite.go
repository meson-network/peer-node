package plugin

import (
	"path/filepath"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/cdn_cache_folder"
	"github.com/meson-network/peer-node/src/common/dbkv"
	"github.com/meson-network/peer-node/src/file_mgr"
)

func initSqlite() error {

	//sqlite_path, sqlite_path_err := storage_mgr.GetInstance().FileExist("db", "peer.db")
	//if sqlite_path_err != nil {
	//	return errors.New("db file not found," + sqlite_path_err.Error() + "\n please check storage folder path , if still not working please re-download your  program and restart")
	//}

	sqlite_path := filepath.Join(cdn_cache_folder.GetInstance().Abs_path, "db", "peer.db")

	err := sqlite_plugin.Init(sqlite_plugin.Config{
		Sqlite_path: sqlite_path,
	}, basic.Logger)
	if err != nil {
		return err
	}

	//auto create table
	err = sqlite_plugin.GetInstance().AutoMigrate(
		&dbkv.DBKVModel{},
		&file_mgr.FileModel{},
	)
	if err != nil {
		return err
	}
	return nil
}
