package plugin

import (
	"os"
	"path/filepath"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/cdn_cache_folder"
	"github.com/meson-network/peer-node/src/common/dbkv"
	"github.com/meson-network/peer-node/src/file_mgr"
)

func initSqlite() error {
	sqlFileFolder := filepath.Join(cdn_cache_folder.GetInstance().Abs_path, "db")
	err := os.MkdirAll(sqlFileFolder, 0777)
	if err != nil {
		return err
	}

	sqlite_path := filepath.Join(sqlFileFolder, "peer.db")

	err = sqlite_plugin.Init(sqlite_plugin.Config{
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
