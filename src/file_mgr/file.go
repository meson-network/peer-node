package file_mgr

import (
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/reference_plugin"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
)

func CreateFile(file *FileModel) (*FileModel, error) {
	if err := sqlite_plugin.GetInstance().Table("file").Create(file).Error; err != nil {
		return nil, err
	}
	GetFile(file.Url_hash, false, true)
	return file, nil
}

func UpdateFile(newData map[string]interface{}, url_hash string) error {
	result := sqlite_plugin.GetInstance().Table("file").Where("url_hash=?", url_hash).Updates(newData)
	if result.Error != nil {
		return result.Error
	}
	GetFile(url_hash, false, true)
	return nil
}

func DeleteFile(url_hash string) error {
	file := &FileModel{Url_hash: url_hash}
	if err := sqlite_plugin.GetInstance().Table("file").Where("url_hash=?", url_hash).Delete(file).Error; err != nil {
		return err
	}
	GetFile(url_hash, false, true)
	return nil
}

func GetFile(url_hash string, fromRef bool, updateRef bool) (*FileModel, error) {

	key := url_hash
	if fromRef {
		basic.Logger.Debugln("GetFile from reference")
		// try to get from reference
		ref_result, _ := reference_plugin.GetInstance().Get(key)
		if ref_result != nil {
			return ref_result.(*FileModel), nil
		}
	}

	//try from db
	var sqlfiles []*FileModel
	basic.Logger.Debugln("GetFile from sqlite")
	sql_err := sqlite_plugin.GetInstance().Table("file").Where("url_hash=?", url_hash).Find(&sqlfiles).Error

	if sql_err != nil {
		basic.Logger.Errorln("GetFile err :", sql_err)
		return nil, sql_err
	} else {
		var result *FileModel = nil
		if len(sqlfiles) != 0 {
			result = sqlfiles[0]
		}
		if updateRef {
			basic.Logger.Debugln("GetFile updateRef")
			reference_plugin.GetInstance().Set(key, result, 1800) //30 mins, long cache time to make things fast
		}
		return result, nil
	}
}
