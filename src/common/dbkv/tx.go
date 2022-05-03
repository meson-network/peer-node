package dbkv

// import (
// 	"context"

// 	"github.com/go-redis/redis/v8"
// 	"github.com/meson-network/peer-node/basic"
// 	"github.com/meson-network/peer-node/plugin/redis_plugin"
// 	"github.com/meson-network/peer-node/plugin/reference_plugin"
// 	"github.com/meson-network/peer-node/tools/smart_cache"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/clause"
// )

// func SetDBKV(tx *gorm.DB, keystr string, value string) error {
// 	err := tx.Table("dbkv").Clauses(clause.OnConflict{UpdateAll: true}).Create(&DBKVModel{Key: keystr, Value: value}).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DeleteDBKV(tx *gorm.DB, keystr string) error {
// 	if err := tx.Table("dbkv").Where(" `key` = ?", keystr).Delete(&DBKVModel{}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func GetDBKV(tx *gorm.DB, keyStr string, fromCache bool, updateCache bool) (*DBKVModel, error) {

// 	//gen_key
// 	ck := smart_cache.NewConnectKey("dbkv")
// 	ck.C_Str(keyStr)

// 	key := redis_plugin.GetInstance().GenKey(ck.String())

// 	if fromCache {
// 		// try to get from reference
// 		result := smart_cache.Ref_Get(reference_plugin.GetInstance(), key)
// 		if result != nil {
// 			basic.Logger.Debugln("GetDBKV hit from reference")
// 			return result.(*DBKVModel), nil
// 		}

// 		redis_result := &DBKVModel{}
// 		// try to get from redis
// 		err := smart_cache.Redis_Get(context.Background(), redis_plugin.GetInstance().ClusterClient, true, key, redis_result)
// 		if err == nil {
// 			basic.Logger.Debugln("GetDBKV hit from redis")
// 			smart_cache.Ref_Set(reference_plugin.GetInstance(), key, redis_result)
// 			return redis_result, nil
// 		} else if err == redis.Nil {
// 			//continue to get from db part
// 		} else if err == smart_cache.TempNil {
// 			return nil, nil
// 		} else {
// 			//redis may broken, just return to keep db safe
// 			return redis_result, err
// 		}
// 	}

// 	//after cache miss ,try from remote database
// 	basic.Logger.Debugln("GetDBKV try from database")

// 	queryResults := []*DBKVModel{}
// 	err := tx.Table("dbkv").Where("`key` = ?", keyStr).Find(&queryResults).Error

// 	if err != nil {
// 		basic.Logger.Errorln("GetDBKV err :", err)
// 		return nil, err
// 	} else {
// 		if len(queryResults) == 0 {
// 			if updateCache {
// 				smart_cache.RR_Set(context.Background(), redis_plugin.GetInstance().ClusterClient, reference_plugin.GetInstance(), true, key, nil, 300)
// 			}
// 			return nil, nil
// 		} else {
// 			if updateCache {
// 				smart_cache.RR_Set(context.Background(), redis_plugin.GetInstance().ClusterClient, reference_plugin.GetInstance(), true, key, queryResults[0], 300)
// 			}
// 			return queryResults[0], nil
// 		}
// 	}
// }
