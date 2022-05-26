package cdn_cache_folder

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/coreservice-io/utils/rand_util"
)

var data = make([]byte, 32*1024, 32*1024) // Initialize an empty byte slice

func (cf *CdnCacheFolder) getFreeSize() int64 {
	cf.sizeLock.RLock()
	defer cf.sizeLock.RUnlock()
	return cf.Cache_provide_size - cf.Cache_used_size
}

func (cf *CdnCacheFolder) addCacheUsedSize(size int64) {
	cf.sizeLock.Lock()
	defer cf.sizeLock.Unlock()
	cf.Cache_used_size += size
}

func (cf *CdnCacheFolder) reduceCacheUsedSize(size int64) {
	cf.sizeLock.Lock()
	defer cf.sizeLock.Unlock()
	cf.Cache_used_size -= size
}

func (cf *CdnCacheFolder) getMesonCacheUsedSize() int64 {
	cf.sizeLock.RLock()
	defer cf.sizeLock.RUnlock()
	return cf.Cache_used_size
}

func (cf *CdnCacheFolder) getSpaceHolderInfo() (size int64, fileNames []string, err error) {
	dirPath := filepath.Join(cf.Abs_path, SpaceHolderFolder)
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return 0, nil, err
	}

	for _, v := range fileInfos {
		if v.IsDir() {
			continue
		}
		size += v.Size()
		fileNames = append(fileNames, filepath.Join(dirPath, v.Name()))
	}

	return size, fileNames, nil
}

func (cf *CdnCacheFolder) CheckFolder(checkLimitGB int) error {
	//check auth
	//make sure folder exist
	spaceHolderFolderPath := filepath.Join(cf.Abs_path, SpaceHolderFolder)
	err := os.MkdirAll(spaceHolderFolderPath, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(spaceHolderFolderPath, 0777)
	if err != nil {
		return err
	}
	fileFolderPath := filepath.Join(cf.Abs_path, CacheFileFolder)
	err = os.MkdirAll(fileFolderPath, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(fileFolderPath, 0777)
	if err != nil {
		return err
	}

	cachedFileSize := getAllFileSizeInFolder(fileFolderPath)
	spaceHolderSize, holdFiles, err := cf.getSpaceHolderInfo()
	if err != nil {
		return err
	}

	//delete space holder file
	defer os.RemoveAll(spaceHolderFolderPath)

	checkLimit := int64(checkLimitGB) * 1024 * 1024 * 1024
	spaceUsed := cachedFileSize + spaceHolderSize
	if spaceUsed >= checkLimit || spaceUsed > cf.Cache_provide_size {
		return nil
	}

	//check 40GB-used
	checkSize := checkLimit - spaceUsed
	freeSize := cf.Cache_provide_size - spaceUsed
	if freeSize < checkSize {
		checkSize = freeSize
	}
	holdFileCount := int(checkSize) / eachHoldFileSize

	fmt.Println("check path:<", cf.Abs_path, ">")
	//print process
	bar := pb.StartNew(holdFileCount)
	bar.SetMaxWidth(100)
	bar.SetWriter(os.Stdout)
	defer bar.Finish()

	for i := 0; i < holdFileCount; i++ {

		fileName, err := genSingleFile(spaceHolderFolderPath)
		if err != nil {
			return err
		}
		if fileName != "" {
			holdFiles = append(holdFiles, fileName)
		}
		bar.Increment()
	}

	return nil
}

func genSingleFile(folderPath string) (string, error) {
	var filePath string
	for {
		fileName := rand_util.GenRandStr(32)
		filePath = filepath.Join(folderPath, fileName)
		info, _ := os.Stat(filePath)
		if info == nil {
			break
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	round := eachHoldFileSize / (32 * 1024)
	for j := 0; j < round; j++ {
		_, err = f.Write(data)
		if err != nil {
			return filePath, err
		}
	}

	return filePath, nil
}

func getAllFileSizeInFolder(folderPath string) int64 {
	fileTotalSize := int64(0)
	filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		f, err := d.Info()
		if err != nil {
			//ignore err
			return nil
		}

		fileTotalSize += f.Size()
		return nil
	})
	return fileTotalSize
}
