package dbkv

//table dbkv
type DBKVModel struct {
	Id    int64  `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
