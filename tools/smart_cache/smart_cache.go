package smart_cache

import (
	"math/rand"

	"github.com/coreservice-io/reference"
)

type temp_nil_error string

func (e temp_nil_error) Error() string { return string(e) }

const local_reference_secs = 5 //don't change this number as 5 is the proper number

// check weather we need do refresh
// the probobility becomes lager when left seconds close to 0
// this goal of this function is to avoid big traffic glitch
func CheckTtlRefresh(secleft int64) bool {
	if secleft > 0 && secleft <= 3 {
		if rand.Intn(int(secleft)*10) == 1 {
			return true
		}
	}
	return false
}

func Ref_Get(localRef *reference.Reference, keystr string) (result interface{}) {
	localvalue, ttl := localRef.Get(keystr)
	if !CheckTtlRefresh(ttl) && localvalue != nil {
		return localvalue
	}
	return nil
}

func Ref_Set(localRef *reference.Reference, keystr string, value interface{}) error {
	return Ref_Set_RTTL(localRef, keystr, value, local_reference_secs)
}

func Ref_Set_RTTL(localRef *reference.Reference, keystr string, value interface{}, ref_ttl_second int64) error {
	return localRef.Set(keystr, value, ref_ttl_second)
}
