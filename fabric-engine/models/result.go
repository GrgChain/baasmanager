package models

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Success = iota
	Fail
)

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

func ResultFail(ctx *gin.Context, err interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": Fail, "data": nil, "msg": err})
}

func ResultFailData(ctx *gin.Context, data interface{}, err interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": Fail, "data": data, "msg": err})
}
