package file_mgr

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/meson-network/peer-node/src/cdn_cache_folder"
)

var deleteEmptyFolderChan = make(chan string, 1000)

func DeleteEmptyFolder(deletedFileAbsPath string) {
	deleteEmptyFolderChan <- deletedFileAbsPath
}

func LoopDeleteEmptyFolder() {
	for {
		folderPath := <-deleteEmptyFolderChan
		checkAndDeleteEmptyFolder(folderPath)
	}
}

func checkAndDeleteEmptyFolder(deletedFileAbsPath string) {
	folder := filepath.Dir(deletedFileAbsPath)
	for {
		if strings.HasSuffix(folder, cdn_cache_folder.CacheFileFolder) {
			return
		}
		dirEntry, err := os.ReadDir(folder)
		if err == nil && len(dirEntry) == 0 {
			err := os.Remove(folder)
			if err != nil {
				return
			}
			folder = filepath.Dir(folder)
		} else {
			return
		}
	}
}
