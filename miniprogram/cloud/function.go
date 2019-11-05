package cloud

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/silenceper/wechat/util"
)

const (
	invokeCloudFunction = "https://api.weixin.qq.com/tcb/invokecloudfunction?access_token=%s&env=%s&name=%s"
)

type FunctionResp struct {
	util.CommonError

	RespData string `json:"resp_data"`
}

/**
 * name: 云函数名称
 * body: 云函数的传入参数
 * 文档地址: https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-http-api/functions/invokeCloudFunction.html
 */
func (cloud *Cloud) InvokeCloudFunction(name string, body interface{}) (res *FunctionResp,e error) {
	if token, err := cloud.GetAccessToken(); err != nil {
		logs.Error("get access token failed. err: %v", err)
		e = err
	} else if resp, err := util.PostJSON(fmt.Sprintf(invokeCloudFunction, token, cloud.Env, name), body); err != nil {
		logs.Error("call invoke cloud function failed. env: %v, name: %v, ", cloud.Env, name)
		e = err
	} else if err = json.Unmarshal(resp, res); err != nil {
		logs.Error("invoke cloud function unmarshal failed. resp: %v, err: %v", string(resp), err)
		e = err
	}
	return res, e
}