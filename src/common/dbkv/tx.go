package dbkv

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetDBKV(tx *gorm.DB, keystr string, value string) error {
	err := tx.Table("dbkv").Clauses(clause.OnConflict{UpdateAll: true}).Create(&DBKVModel{Key: keystr, Value: value}).Error
	if err != nil {
		return err
	}
	GetDBKV(tx, keystr, true)
	return nil
}

func DeleteDBKV(tx *gorm.DB, keystr string) error {
	if err := tx.Table("dbkv").Where(" `key` = ?", keystr).Delete(&DBKVModel{}).Error; err != nil {
		return err
	}
	GetDBKV(tx, keystr, true)
	return nil
}

func GetDBKV(tx *gorm.DB, keyStr string, updateCache bool) (*DBKVModel, error) {

	key := "dbkv" + ":" + keyStr

	result, _ := reference_plugin.GetInstance().Get(key)
	if result != nil {
		basic.Logger.Debugln("GetDBKV hit from reference")
		return result.(*DBKVModel), nil
	}

	//after cache miss ,try from sqlite
	basic.Logger.Debugln("GetDBKV try from database")

	queryResults := []*DBKVModel{}
	err := tx.Table("dbkv").Where("`key` = ?", keyStr).Find(&queryResults).Error

	if err != nil {
		basic.Logger.Errorln("GetDBKV err :", err)
		return nil, err
	} else {
		var result *DBKVModel = nil
		if len(queryResults) != 0 {
			result = queryResults[0]
		}
		reference_plugin.GetInstance().Set(key, result, 900)
		return result, nil
	}
}
