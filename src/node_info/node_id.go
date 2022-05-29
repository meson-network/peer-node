package node_info

import (
	"github.com/coreservice-io/utils/rand_util"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/common/dbkv"
)

var node_id string

func GetNodeId() string {
	return node_id
}

func InitNode() error {
	db_node_id, _ := dbkv.GetKey(sqlite_plugin.GetInstance(), "node_id", false, false)

	if db_node_id == "" {
		//create a node_id in dbkv
		node_id = rand_util.GenRandStr(16)
		db_err := dbkv.SetDBKV(sqlite_plugin.GetInstance(), "node_id", node_id)
		if db_err != nil {
			return db_err
		}
	} else {
		node_id = db_node_id
	}
	return nil
}
