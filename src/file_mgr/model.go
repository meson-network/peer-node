package file_mgr

const STATUS_DOWNLOADED = "DOWNLOADED"
const STATUS_DOWNLOADING = "DOWNLOADING"

const TYPE_PUBLIC = "PUBLIC"
const TYPE_PRIVATE = "PRIVATE"

type FileModel struct {
	Url_hash               string `gorm:"url_hash"`
	Last_req_unixtime      int64  `gorm:"last_req_unixtime"`
	Last_scan_unixtime     int64  `gorm:"last_scan_unixtime"`
	Last_download_unixtime int64  `gorm:"last_download_unixtime"`
	Size_byte              int64  `gorm:"size_byte"`
	Rel_path               string `gorm:"rel_path"`
	Status                 string `gorm:"status"`
	Type                   string `gorm:"type"`
}
