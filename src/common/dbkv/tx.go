package dbkv

import (
	"errors"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"github.com/meson-network/peer-node/tools/smart_cache"
	"gorm.io/gorm"
)

func SetDBKV(tx *gorm.DB, keystr string, value string) error {
	result, err := QueryDBKV(tx, nil, &[]string{keystr}, false, false)
	if err == nil && len(result.Kv) > 0 {
		//exist update
		err := tx.Table("dbkv").Where("`key`=?", keystr).Updates(map[string]interface{}{"value": value}).Error
		if err != nil {
			basic.Logger.Errorln("SetDBKV update error:", err)
			return err
		}
		//refresh
		GetKey(tx, keystr, false, true)
		return nil
	} else {
		//create
		err := tx.Table("dbkv").Create(&DBKVModel{Key: keystr, Value: value}).Error
		if err != nil {
			basic.Logger.Errorln("SetDBKV create error:", err)
			return err
		}
		//refresh
		GetKey(tx, keystr, false, true)
		return nil
	}
}

func DeleteDBKV(tx *gorm.DB, keystr string) error {
	if err := tx.Table("dbkv").Where(" `key` = ?", keystr).Delete(&DBKVModel{}).Error; err != nil {
		return err
	}
	GetKey(tx, keystr, false, true)
	return nil
}

func GetKey(tx *gorm.DB, key string, fromCache bool, updateCache bool) (string, error) {
	result, err := QueryDBKV(tx, nil, &[]string{key}, fromCache, updateCache)
	if err != nil {
		return "", err
	}
	if len(result.Kv) == 0 {
		return "", errors.New("key not exist")
	}
	return result.Kv[0].Value, nil
}

type QueryKvResult struct {
	Kv         []*DBKVModel
	TotalCount int64
}

func QueryDBKV(tx *gorm.DB, id *int64, keys *[]string, fromCache bool, updateCache bool) (*QueryKvResult, error) {

	//gen_key
	ck := smart_cache.NewConnectKey("dbkv")
	ck.C_Int64_Ptr("id", id).
		C_Str_Array_Ptr("keys", keys)

	key := ck.String()

	if fromCache {
		// try to get from reference_plugin
		result := smart_cache.Ref_Get(reference_plugin.GetInstance(), key)
		if result != nil {
			basic.Logger.Debugln("QueryUser hit from reference_plugin")
			return result.(*QueryKvResult), nil
		}
	}

	//after cache miss ,try from remote database
	basic.Logger.Debugln("QueryDBKV try from database")

	queryResult := &QueryKvResult{
		Kv:         []*DBKVModel{},
		TotalCount: 0,
	}

	query := tx.Table("dbkv")
	if id != nil {
		query.Where("id = ?", *id)
	}
	if keys != nil {
		query.Where("dbkv.key IN ?", *keys)
	}

	query.Count(&queryResult.TotalCount)

	err := query.Find(&queryResult.Kv).Error
	if err != nil {
		basic.Logger.Errorln("QueryDBKV err :", err)
		return nil, err
	} else {
		if updateCache {
			smart_cache.Ref_Set(reference_plugin.GetInstance(), key, queryResult)
		}
		return queryResult, nil
	}
}
