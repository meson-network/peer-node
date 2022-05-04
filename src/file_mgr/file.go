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
	GetFile(file.Hash, false, true)
	return file, nil
}

func UpdateFile(newData map[string]interface{}, hash string) error {
	result := sqlite_plugin.GetInstance().Table("file").Where("hash=?", hash).Updates(newData)
	if result.Error != nil {
		return result.Error
	}
	GetFile(hash, false, true)
	return nil
}

func DeleteFile(hash string) error {
	file := &FileModel{Hash: hash}
	if err := sqlite_plugin.GetInstance().Table("file").Where("hash=?", hash).Delete(file).Error; err != nil {
		return err
	}
	GetFile(hash, false, true)
	return nil
}

func GetFile(hash string, fromRef bool, updateRef bool) (*FileModel, error) {

	key := hash

	if fromRef {
		basic.Logger.Debugln("GetFile from reference")
		// try to get from reference
		ref_result, _ := reference_plugin.GetInstance().Get(key)
		if ref_result != nil {
			return ref_result.(*FileModel), nil
		}
	}

	var sqlfiles []*FileModel

	basic.Logger.Debugln("GetFile from sqlite")
	//try from db
	sql_err := sqlite_plugin.GetInstance().Table("file").Where("hash=?", hash).Find(&sqlfiles).Error

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
			reference_plugin.GetInstance().Set(key, result, 300)
		}
		return result, nil
	}
}
