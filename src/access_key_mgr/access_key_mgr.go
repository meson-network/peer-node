package access_key_mgr

import (
	"fmt"

	"github.com/coreservice-io/utils/rand_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/meson-network/peer-node/plugin/sqlite_plugin"
	"github.com/meson-network/peer-node/src/common/dbkv"
)

type AccessKeyMgr struct {
	currentKey  string
	previousKey string
}

var instanceMap = map[string]*AccessKeyMgr{}

func GetInstance() *AccessKeyMgr {
	return instanceMap["default"]
}

func GetInstance_(name string) *AccessKeyMgr {
	return instanceMap[name]
}

func Init(curKey string, preKey string) error {
	return Init_("default", curKey, preKey)
}

func Init_(name string, curKey string, preKey string) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("accessKeyMgr instance <%s> has already initialized", name)
	}

	instanceMap[name] = &AccessKeyMgr{
		currentKey:  curKey,
		previousKey: preKey,
	}
	return nil
}

//GenNewRandomKey generate a new random key as currentKey, and the old one will be record as previousKey
func (r *AccessKeyMgr) GenNewRandomKey() string {
	randKey := rand_util.GenRandStr(10)
	if r.currentKey == "" {
		r.previousKey = randKey
	} else {
		r.previousKey = r.currentKey
	}
	r.currentKey = randKey
	err := dbkv.SetDBKV(sqlite_plugin.GetInstance(), "access_key", r.currentKey+","+r.previousKey)
	if err != nil {
		basic.Logger.Errorln("GenNewRandomKey dbkv.SetDBKV error:", err)
	}
	return randKey
}

//CheckRandomKey check the given randomKey is equal to the current key or previous key
func (r *AccessKeyMgr) CheckRandomKey(inputKey string) bool {
	if inputKey == r.currentKey || inputKey == r.previousKey {
		return true
	}
	return false
}

func (r *AccessKeyMgr) GetRandomKey() (currentKey string, previousKey string) {
	if r.currentKey == "" {
		r.GenNewRandomKey()
	}
	return r.currentKey, r.previousKey
}
