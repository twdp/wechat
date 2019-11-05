package cloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/silenceper/wechat/util"
)

const (
	databaseImportUrl        = "https://api.weixin.qq.com/tcb/databasemigrateimport?access_token=%s"
	databaseExportUrl        = "https://api.weixin.qq.com/tcb/databasemigrateexport?access_token=%s"
	databaseMigrateQueryInfo = "https://api.weixin.qq.com/tcb/databasemigratequeryinfo?access_token=%s"
)

type FileType int8
type ConflictMode int8

const (
	// file_type
	Json FileType = 1
	Csv  FileType = 2

	// conflict_mode
	Insert ConflictMode = 1
	Upsert ConflictMode = 2
)

// 导入参数
type ImportDataParam struct {
	// 导入collection名
	CollectionName string `json:"collection_name"`

	// 导入文件路径(导入文件需先上传到同环境的存储中，可使用开发者工具或 HTTP API的上传文件 API上传）
	FilePath string `json:"file_path"`

	// 导入文件类型，文件格式参考数据库导入指引中的文件格式部分
	FileType FileType `json:"file_type"`

	// 是否在遇到错误时停止导入
	StopOnError bool `json:"stop_on_error"`

	// 冲突处理模式
	ConflictMode ConflictMode `json:"conflict_mode"`
}

/**
 * 数据库导入
 */
func (cloud *Cloud) DatabaseImport(importDataParam ImportDataParam) (jobId int64, err error) {
	param := struct {
		ImportDataParam
		Env string `json:"env"`
	}{
		ImportDataParam: importDataParam,
		Env:             cloud.Env,
	}

	res := struct {
		util.CommonError
		JobId int64 `json:"job_id"`
	}{}

	if token, e := cloud.GetAccessToken(); e != nil {
		return 0, e
	} else if resp, e := util.PostJSON(fmt.Sprintf(databaseImportUrl, token), param); e != nil {
		logs.Error("import database failed. param: %v, err: %v", param, e)
		err = e
	} else if e := json.Unmarshal(resp, &res); e != nil {
		logs.Error("import database unmarshal failed. param: %v, res: %v, e: %v", param, string(resp), e)
		err = e
	} else if fail, _, msg := util.IsError(res.CommonError); fail {
		logs.Error("import return err. res: %v", res)

		err = errors.New(msg)
	} else {
		jobId = res.JobId
	}
	return jobId, err
}

// 导出参数
type ExportDataParam struct {
	// 导出文件路径（文件会导出到同环境的云存储中，可使用获取下载链接 API 获取下载链接）
	FilePath string `json:"file_path"`

	// 导出文件类型，文件格式参考数据库导入指引中的文件格式部分
	FileType FileType `json:"file_type"`

	// 导出条件 查询语句语法与数据库 API相同
	// api: https://developers.weixin.qq.com/miniprogram/dev/wxcloud/reference-client-api/database/#%E6%95%B0%E6%8D%AE%E5%BA%93
	Query string `json:"query"`
}

/**
 * 数据库导出
 */
func (cloud *Cloud) DatabaseExport(exportDataParam ExportDataParam) (jobId int64, err error) {
	param := struct {
		ExportDataParam
		Env string `json:"env"`
	}{
		ExportDataParam: exportDataParam,
		Env:             cloud.Env,
	}

	res := struct {
		util.CommonError
		JobId int64 `json:"job_id"`
	}{}

	if token, e := cloud.GetAccessToken(); e != nil {
		return 0, e
	} else if resp, e := util.PostJSON(fmt.Sprintf(databaseExportUrl, token), param); e != nil {
		logs.Error("export database failed. param: %v, err: %v", param, e)
		err = e
	} else if e := json.Unmarshal(resp, &res); e != nil {
		logs.Error("export database unmarshal failed. param: %v, res: %v, e: %v", param, string(resp), e)
		err = e
	} else if fail, _, msg := util.IsError(res.CommonError); fail {
		logs.Error("export return err. res: %v", res)
		err = errors.New(msg)
	} else {
		jobId = res.JobId
	}
	return
}
