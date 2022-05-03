package errors

import (
	"errors"
	"strconv"
)

// status must be [ -10000 to -99999 ] ,if status wrong will be set to -10000 as default
func NewStatusError(status int64, msg string) error {
	if status < -9999 && status > -100000 {

	} else {
		status = -10000
	}
	return errors.New(strconv.FormatInt(status, 10) + "," + msg)
}

//if err is nil then return 0,nil
//if not an StatusError then return -1,error
//e.g "-10000,error found" => -10000, err("error found")
//e.g "this is error" => -1, err("this is error")
//e.g "101,this is positive err" => -1,error("101,this is positive err")
//e.g "-10,this is not status error" => -1,error("-10,this is not status error")
func ResolveStatusError(err error) (int64, error) {
	if err == nil {
		return 0, nil
	}

	errMsg := err.Error()
	if len(errMsg) < 7 || errMsg[0:1] != "-" || errMsg[6:7] != "," {
		return -1, err
	}

	status, parse_err := strconv.ParseInt(errMsg[0:6], 10, 64)
	if parse_err != nil {
		return -1, err
	}

	return status, errors.New(errMsg[7:])
}
