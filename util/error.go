package util

import (
	"encoding/json"
	"fmt"
)

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// DecodeWithCommonError 将返回值按照CommonError解析
func DecodeWithCommonError(response []byte, apiName string) (err error) {
	var commError CommonError
	err = json.Unmarshal(response, &commError)
	if err != nil {
		return
	}
	if commError.ErrCode != 0 {
		return fmt.Errorf("%s Error , errcode=%d , errmsg=%s", apiName, commError.ErrCode, commError.ErrMsg)
	}
	return nil
}

/**
 * 判断是否出错
 * result true: 出错, false: 为出错
 * code 返回的code
 * msg 返回的错误信息
 */
func IsError(commonError CommonError) (result bool, code int64, msg string) {
	if commonError.ErrCode == 0 {
		result = false
	} else {
		result = true
		code = commonError.ErrCode
		msg = commonError.ErrMsg
	}
	return result, code, msg
}
