package dbkv

type DBKVModel struct {
	Id    int64  `gorm:"primarykey"`
	Key   string `gorm:"index;unique"`
	Value string
}
