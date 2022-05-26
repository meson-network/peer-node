package dbkv

//table dbkv
type DBKVModel struct {
	Id    int64  `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"index;unique"`
	Value string `json:"value"`
}

func (DBKVModel) TableName() string {
	return "dbkv"
}
