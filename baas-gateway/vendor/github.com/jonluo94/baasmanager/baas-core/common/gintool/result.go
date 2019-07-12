package gintool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Success = iota
	Fail
)

type RespData struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
}

func ResultMap(ctx *gin.Context, m map[string]interface{}) {
	ctx.JSON(http.StatusOK, m)
}

func ResultMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{"code": Success, "data": nil, "msg": msg})
}
func ResultOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": Success, "data": data, "msg": ""})
}
func ResultOkMsg(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, gin.H{"code": Success, "data": data, "msg": msg})
}

/**
  @param  total  :总数
  @param  pageNum : 当前页
*/
type RespDataList struct {
	RespData
	Total int64 `json:"total"`
}

func ResultList(ctx *gin.Context, data interface{}, total int64) {
	ctx.JSON(http.StatusOK, gin.H{"code": Success, "data": data, "msg": "", "total": total})
}

type RespDataPager struct {
	RespData
	Pager interface{} `json:"pager"`
}

func ResultPageList(ctx *gin.Context, data interface{}, pager interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": Success, "data": data, "msg": "", "pager": pager})
}

func ResultFail(ctx *gin.Context, err interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": Fail, "data": nil, "msg": err})
}

func ResultFailData(ctx *gin.Context, data interface{}, err interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": Fail, "data": data, "msg": err})
}
