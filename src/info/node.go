package info

import (
	"github.com/coreservice-io/utils/rand_util"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/common/dbkv"
	"github.com/meson-network/peer_common/info"
)

var node_id string

func GetNodeId() string {
	return node_id
}

func InitNode() error {
	db_node_id, db_node_id_err := dbkv.GetDBKV(sqlite_plugin.GetInstance(), info.NODE_ID, true)
	if db_node_id_err != nil {
		return db_node_id_err
	}
	if db_node_id == nil {
		//create a node_id in dbkv
		node_id = rand_util.GenRandStr(16)
		db_err := dbkv.SetDBKV(sqlite_plugin.GetInstance(), info.NODE_ID, node_id)
		if db_err != nil {
			return db_err
		}
	} else {
		node_id = db_node_id.Value
	}
	return nil
}
