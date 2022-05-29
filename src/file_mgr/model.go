package file_mgr

const STATUS_DOWNLOADED = "DOWNLOADED"
const STATUS_DOWNLOADING = "DOWNLOADING"

//const TYPE_PUBLIC = "PUBLIC"
//const TYPE_PRIVATE = "PRIVATE"

type FileModel struct {
	File_hash         string `json:"file_hash" gorm:"primaryKey"`
	Last_req_unixtime int64  `json:"last_req_unixtime" gorm:"index"`
	//Last_scan_unixtime     int64  `json:"last_scan_unixtime" gorm:"index"`
	//Last_download_unixtime int64  `json:"last_download_unixtime"`
	Size_byte int64  `json:"size_byte" gorm:"index"`
	Rel_path  string `json:"rel_path"`
	Status    string `json:"status" gorm:"index"`
	//Type                   string `json:"type"`
}

func (FileModel) TableName() string {
	return "file"
}
