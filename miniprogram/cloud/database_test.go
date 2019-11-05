package cloud

import (
	"fmt"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"testing"
)

var c = cache.NewMemory()

func TestCloud_DatabaseExport(t *testing.T) {

	//配置微信参数
	config := &wechat.Config{
		AppID:          "wxd7039a6dfc9358ce",
		AppSecret:      "9b78b7e40c414ebda37bba8c353244d9",
		Cache:          c,
	}
	wc := wechat.NewWechat(config)

	cloud := &Cloud{
		Context: wc.Context,
		Env:     "tw-001",
	}

	fmt.Println(cloud.DatabaseExport(ExportDataParam{
		FilePath: "test",
		FileType: Json,
		Query:    `db.collection('question').get()`,
	}))
}

func TestCloud_DatabaseMigrateQueryInfo(t *testing.T) {
	//配置微信参数
	config := &wechat.Config{
		AppID:          "wxd7039a6dfc9358ce",
		AppSecret:      "9b78b7e40c414ebda37bba8c353244d9",
		Cache:          c,
	}
	wc := wechat.NewWechat(config)

	cloud := &Cloud{
		Context: wc.Context,
		Env:     "tw-001",
	}

	fmt.Println(cloud.DatabaseMigrateQueryInfo(100164445))
}