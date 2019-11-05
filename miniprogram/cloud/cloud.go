package cloud

import "github.com/silenceper/wechat/context"

// MiniProgram struct extends context
type Cloud struct {
	*context.Context
	Env string // 云开发 env 环境id
}

// NewMiniProgram 实例化小程序接口
func NewCloud(context *context.Context) *Cloud {
	miniProgram := new(Cloud)
	miniProgram.Context = context
	return miniProgram
}