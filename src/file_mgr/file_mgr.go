package file_mgr

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/src/cdn_cache_folder"
)

func RemoveFileFromDisk(fileHash string) {
	fileAbsPath := GetFileAbsPath(fileHash)
	os.Remove(fileAbsPath)
	os.Remove(fileAbsPath + ".header")
}

func GetFileAbsPath(file_hash string) string {
	return filepath.Join(cdn_cache_folder.GetInstance().Abs_path, UrlHashToPublicFileRelPath(file_hash), file_hash)
}

func CleanDownloadingFiles() error {
	result, err := QueryFile(nil, nil, &[]string{STATUS_DOWNLOADING}, nil, 0, 0, false, false)
	if err != nil {
		return err
	}

	if len(result.Files) == 0 {
		return nil
	}

	for _, v := range result.Files {
		//todo get abs path
		fileAbsPath := GetFileAbsPath(v.File_hash)

		os.Remove(fileAbsPath)
		os.Remove(fileAbsPath + ".header")

		DeleteFile(v.File_hash)
		DeleteEmptyFolder(fileAbsPath)
	}

	return nil

}

type DiskFile struct {
	absPath  string
	fileName string
}

//ScanLeakFiles Scan disk and delete files not in db
func ScanLeakFiles() {

	files := []*DiskFile{}
	rootPath := filepath.Join(cdn_cache_folder.GetInstance().Abs_path, cdn_cache_folder.CacheFileFolder)
	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		fileName := d.Name()
		fileExtension := filepath.Ext(fileName)
		//extension is not ".header"
		if fileExtension != "" && fileExtension != ".header" {
			//not cache file
			err := os.Remove(path)
			if err != nil {

				basic.Logger.Errorln("scanLeakFiles remove file error:", err, "path:", path)

			} else {
				DeleteEmptyFolder(path)
			}

			return nil
		}
		//fileBaseName is not 16 hash
		fileBase := strings.TrimSuffix(fileName, ".header")
		if len(fileBase) != 16 {
			err := os.Remove(path)
			if err != nil {

				basic.Logger.Errorln("scanLeakFiles remove file error:", err, "path:", path)

			} else {
				DeleteEmptyFolder(path)
			}
			return nil
		}

		if strings.HasSuffix(fileName, ".header") {
			return nil
		}

		diskFile := &DiskFile{
			absPath:  path,
			fileName: d.Name(),
		}

		files = append(files, diskFile)
		if len(files) >= 100 {
			checkFileLeak(files)
			files = []*DiskFile{}
			time.Sleep(10 * time.Second)
		}

		return nil
	})

	if len(files) > 0 {
		checkFileLeak(files)
		files = []*DiskFile{}
		time.Sleep(10 * time.Second)
	}

}

func checkFileLeak(files []*DiskFile) {
	fileNames := []string{}
	for _, file := range files {
		fileNames = append(fileNames, file.fileName)
	}

	result, err := QueryFile(nil, nil, nil, &fileNames, 0, 0, false, false)
	if err != nil {
		basic.Logger.Errorln("checkFileLeak QueryFile err:", err)
	}

	//all file find in db
	if len(fileNames) == len(result.Files) {
		return
	}

	nameHashMap := map[string]struct{}{}
	for _, v := range result.Files {
		nameHashMap[v.File_hash] = struct{}{}
	}

	for _, file := range files {
		_, exist := nameHashMap[file.fileName]
		if !exist {
			os.Remove(file.absPath)
			os.Remove(file.absPath + ".header")

			basic.Logger.Debugln("deleted leak file", file.absPath)

			DeleteEmptyFolder(file.absPath)
		}
	}
}
