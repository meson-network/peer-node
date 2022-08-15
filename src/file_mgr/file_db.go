package file_mgr

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/tools/smart_cache"
)

func CreateFile(file *FileModel) (*FileModel, error) {
	if err := sqlite_plugin.GetInstance().Table("file").Create(file).Error; err != nil {
		return nil, err
	}
	GetFile(file.File_hash, false, true)
	return file, nil
}

func UpdateFile(newData map[string]interface{}, file_hash string) error {
	result := sqlite_plugin.GetInstance().Table("file").Where("file_hash=?", file_hash).Updates(newData)
	if result.Error != nil {
		return result.Error
	}
	GetFile(file_hash, false, true)
	return nil
}

func DeleteFile(file_hash string) error {
	file := &FileModel{File_hash: file_hash}
	if err := sqlite_plugin.GetInstance().Table("file").Where("file_hash=?", file_hash).Delete(file).Error; err != nil {
		return err
	}
	GetFile(file_hash, false, true)
	return nil
}

func GetFile(file_hash string, fromRef bool, updateRef bool) (*FileModel, error) {

	result, err := QueryFile(&file_hash, nil, nil, nil, 1, 0, fromRef, updateRef)
	if err != nil {
		return nil, err
	}

	if len(result.Files) == 0 {
		return nil, nil
	}

	return result.Files[0], nil

}

type QueryFileResult struct {
	Files      []*FileModel
	TotalCount int64
}

func QueryFile(file_hash *string, less_than_req_unixtime *int64, status *[]string, file_hashs *[]string, limit int, offset int, fromRef bool, updateRef bool) (*QueryFileResult, error) {
	//gen_key
	ck := smart_cache.NewConnectKey("cachedFile")
	ck.C_Str_Ptr("file_hash", file_hash).
		C_Int64_Ptr("less_than_req_unixtime", less_than_req_unixtime).
		C_Str_Array_Ptr("status", status).
		C_Str_Array_Ptr("file_hashs", file_hashs).
		C_Int(limit).
		C_Int(offset)
	key := ck.String()

	if fromRef {
		basic.Logger.Debugln("GetFile from reference")
		// try to get from reference
		ref_result, _ := reference_plugin.GetInstance().Get(key)
		if ref_result != nil {
			return ref_result.(*QueryFileResult), nil
		}
	}

	//try from db
	queryResult := &QueryFileResult{
		Files:      []*FileModel{},
		TotalCount: 0,
	}
	basic.Logger.Debugln("GetFile from sqlite")

	query := sqlite_plugin.GetInstance().Table("file")

	if file_hash != nil {
		query.Where("file_hash = ?", *file_hash)
	}

	if less_than_req_unixtime != nil {
		query.Where("last_req_unixtime < ?", *less_than_req_unixtime)
	}

	if status != nil {
		query.Where("status IN ?", *status)
	}

	if file_hashs != nil {
		query.Where("file_hash IN ?", *file_hashs)
	}

	query.Count(&queryResult.TotalCount)
	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}

	err := query.Find(&queryResult.Files).Error

	if err != nil {
		basic.Logger.Errorln("GetFile err :", err)
		return nil, err
	} else {
		if updateRef {
			basic.Logger.Debugln("GetFile updateRef")
			reference_plugin.GetInstance().Set(key, queryResult, 1800) //30 mins, long cache time to make things fast //todo 5sec for test
		}
		return queryResult, nil
	}
}

func QueryExpireFile(timeline_sec int64, limit int, offset int) (*QueryFileResult, error) {
	//try from db
	queryResult := &QueryFileResult{
		Files:      []*FileModel{},
		TotalCount: 0,
	}
	basic.Logger.Debugln("GetFile from sqlite")

	query := sqlite_plugin.GetInstance().Table("file")
	query.Where("status IN ?", []string{STATUS_DOWNLOADED})
	query.Where("(last_req_unixtime+no_access_maintain_sec) < ?", timeline_sec)

	query.Count(&queryResult.TotalCount)
	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}

	err := query.Find(&queryResult.Files).Error

	if err != nil {
		basic.Logger.Errorln("GetFile err :", err)
		return nil, err
	} else {
		return queryResult, nil
	}
}
