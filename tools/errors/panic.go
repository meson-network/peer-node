package errors

import (
	"crypto/md5"
	"encoding/hex"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coreservice-io/job"
	"github.com/meson-network/peer-node/basic"
)

var errMap sync.Map

//to do panic_err handler :default function
func PanicHandler(panic_err interface{}) {
	basic.Logger.Errorln(panic_err)

	//record panic
	var errStr string
	switch e := panic_err.(type) {
	case string:
		errStr = e
	case runtime.Error:
		errStr = e.Error()
	case error:
		errStr = e.Error()
	default:
		errStr = "recovered (default) panic"
	}

	_ = errStr
	//recordPanicStack(errStr, string(debug.Stack()))
}

func ScheduleUploadPanic() {
	job.Start(
		//job process
		"uploadPanic",
		scheduleUploadPanic,
		//onPanic callback
		PanicHandler,
		300,
		// job type
		// job.TYPE_PANIC_REDO  auto restart if panic
		// job.TYPE_PANIC_RETURN  stop if panic
		job.TYPE_PANIC_REDO,
		// check continue callback, the job will stop running if return false
		// the job will keep running if this callback is nil
		nil,
		// onFinish callback
		nil,
	)
}

func scheduleUploadPanic() {
	errs := [][]string{}

	errMap.Range(func(key, value interface{}) bool {
		//get err
		v, ok := value.([]string)
		if ok {
			errs = append(errs, v)
		}
		//delete
		errMap.Delete(key)
		return true
	})

	if len(errs) == 0 {
		return
	}

	//upload to es or server

}

func recordPanicStack(panicstr string, stack string) {

	errors := []string{panicstr}
	errstr := panicstr

	errors = append(errors, "last err unix-time:"+strconv.FormatInt(time.Now().Unix(), 10))

	lines := strings.Split(stack, "\n")
	maxlines := len(lines)
	if maxlines >= 100 {
		maxlines = 100
	}

	if maxlines >= 3 {
		for i := 2; i < maxlines; i = i + 2 {
			fomatstr := strings.ReplaceAll(lines[i], "	", "")
			errstr = errstr + "#" + fomatstr
			errors = append(errors, fomatstr)
		}
	}

	h := md5.New()
	h.Write([]byte(errstr))
	errhash := hex.EncodeToString(h.Sum(nil))
	_ = errhash

	errMap.Store(errhash, errors)
}
