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
	databaseCollectionAdd    = "https://api.weixin.qq.com/tcb/databasecollectionadd?access_token=%s"
	databaseCollectionDelete = "https://api.weixin.qq.com/tcb/databasecollectiondelete?access_token=%s"
	databaseCollectionGet    = "https://api.weixin.qq.com/tcb/databasecollectionget?access_token=%s"
	databaseAdd              = "https://api.weixin.qq.com/tcb/databaseadd?access_token=%s"
	databaseDelete           = "https://api.weixin.qq.com/tcb/databasedelete?access_token=%s"
	databaseUpdate           = "https://api.weixin.qq.com/tcb/databaseupdate?access_token=%s"
	databaseQuery            = "https://api.weixin.qq.com/tcb/databasequery?access_token=%s"
	databaseCount            = "https://api.weixin.qq.com/tcb/databasecount?access_token=%s"
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

type MigrateQueryInfo struct {
	util.CommonError
	Status        string `json:"status"`
	RecordSuccess int64  `json:"record_success"`
	RecordFail    int64  `json:"record_fail"`
	ErrMsg        string `json:"err_msg"`
	FileUrl       string `json:"file_url"`
}

/**
 * 查询导出状态
 */
func (cloud *Cloud) DatabaseMigrateQueryInfo(jobId int64) (m *MigrateQueryInfo, err error) {
	m = &MigrateQueryInfo{}

	if token, e := cloud.GetAccessToken(); e != nil {
		return m, e
	} else if resp, e := util.PostJSON(fmt.Sprintf(databaseMigrateQueryInfo, token), struct {
		Env   string `json:"env"`
		JobId int64  `json:"job_id"`
	}{
		Env:   cloud.Env,
		JobId: jobId,
	}); e != nil {
		logs.Error("query migrate database failed. param: %v, err: %v", jobId, e)
		err = e
	} else if e := json.Unmarshal(resp, m); e != nil {
		logs.Error("query migrate database unmarshal failed. param: %v, res: %v, e: %v", jobId, string(resp), e)
		err = e
	} else if fail, _, msg := util.IsError(m.CommonError); fail {
		logs.Error("query migrate return err. res: %v", m)
		err = errors.New(msg)
	}
	return
}

/**
 * 添加集合
 */
func (cloud *Cloud) DatabaseCollectionAdd(name string) error {
	if token, e := cloud.GetAccessToken(); e != nil {
		return e
	} else if resp, e := util.PostJSON(fmt.Sprintf(databaseCollectionAdd, token), struct {
		Env            string `json:"env"`
		CollectionName string `json:"collection_name"`
	}{
		Env:            cloud.Env,
		CollectionName: name,
	}); e != nil {
		logs.Error("collection add failed. param: %v, err: %v", name, e)
		return e
	} else if err := util.DecodeWithCommonError(resp, "添加集合"); err != nil {
		logs.Error("add collection failed. resp: %s", string(resp))
		return err
	}
	return nil
}

/**
 * 删除集合
 */
func (cloud *Cloud) DatabaseCollectionDelete(name string) error {
	if token, e := cloud.GetAccessToken(); e != nil {
		return e
	} else if resp, e := util.PostJSON(fmt.Sprintf(databaseCollectionDelete, token), struct {
		Env            string `json:"env"`
		CollectionName string `json:"collection_name"`
	}{
		Env:            cloud.Env,
		CollectionName: name,
	}); e != nil {
		logs.Error("collection delete failed. param: %v, err: %v", name, e)
		return e
	} else if err := util.DecodeWithCommonError(resp, "删除集合"); err != nil {
		logs.Error("delete collection failed. resp: %s", string(resp))
		return err
	}
	return nil
}

type Pager struct {
	Offset int64 `json:"Offset"`
	Limit  int64 `json:"Limit"`
	Total  int64 `json:"Total"`
}

type Collection struct {
	Name       string `json:"name"`
	Count      int64  `json:"count"`
	Size       int64  `json:"size"`
	IndexCount int64  `json:"index_count"`
	IndexSize  int64  `json:"index_size"`
}

type CollectionGetResp struct {
	util.CommonError

	Pager Pager `json:"pager"`

	Collections []Collection `json:"collections"`
}

/**
 * 获取特定云环境下集合信息
 */
func (cloud *Cloud) DatabaseCollectionGet(limit, offset int64) (*CollectionGetResp, error) {
	if token, e := cloud.GetAccessToken(); e != nil {
		return nil, e
	} else if resp, e := util.PostJSON(fmt.Sprintf(databaseCollectionGet, token), struct {
		Env    string `json:"env"`
		Limit  int64  `json:"limit"`
		Offset int64  `json:"offset"`
	}{
		Env:    cloud.Env,
		Limit:  limit,
		Offset: offset,
	}); e != nil {
		logs.Error("collection get failed. limit: %d, offset: %d, err: %v", limit, offset, e)
		return nil, e
	} else {
		result := &CollectionGetResp{}
		if err := json.Unmarshal(resp, result); err != nil {
			logs.Error("unmarshal collection get resp failed. resp: %v", string(resp))
			return nil, errors.New("获取失败")
		} else if result.ErrCode != 0 {
			logs.Error(" collection get resp failed. resp: %v", string(resp))
		}
		return result, errors.New(result.ErrMsg)
	}

}
