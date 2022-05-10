package plugin

import (
	"errors"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/storage_mgr"
)

func initSqlite() error {

	sqlite_path, sqlite_path_err := storage_mgr.CheckPath("db", "peer.db")
	if sqlite_path_err != nil {
		return errors.New("db file not found," + sqlite_path_err.Error() + "\n please check storage folder path , if still not working please re-download your  program and restart")
	}

	return sqlite_plugin.Init(sqlite_plugin.Config{
		Sqlite_path: sqlite_path,
	}, basic.Logger)
}
