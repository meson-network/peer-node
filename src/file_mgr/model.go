package file_mgr

const STATUS_DOWNLOADED = "DOWNLOADED"
const STATUS_DOWNLOADING = "DOWNLOADING"

//const TYPE_PUBLIC = "PUBLIC"
//const TYPE_PRIVATE = "PRIVATE"

type FileModel struct {
	File_hash              string `json:"file_hash" gorm:"primaryKey"`
	Last_req_unixtime      int64  `json:"last_req_unixtime" gorm:"index"`
	No_access_maintain_sec int64  `json:"no_access_maintain_sec" gorm:"index"`
	Size_byte              int64  `json:"size_byte" gorm:"index"`
	Rel_path               string `json:"rel_path"`
	Status                 string `json:"status" gorm:"index"`
}

func (FileModel) TableName() string {
	return "file"
}
