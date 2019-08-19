package controller

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
)

func (a *ApiController) RoleList(ctx *gin.Context) {

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		gintool.ResultFail(ctx, "page error")
		return
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		gintool.ResultFail(ctx, "limit error")
		return
	}
	name := ctx.Query("name")

	b, list, total := a.roleService.GetList(&entity.Role{Name: name}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) RoleAllList(ctx *gin.Context) {

	b, list := a.roleService.GetAll()
	if b {
		gintool.ResultOk(ctx, list)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) RoleAdd(ctx *gin.Context) {

	role := new(entity.Role)

	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.roleService.Add(role)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) RoleUpdate(ctx *gin.Context) {

	role := new(entity.Role)

	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.roleService.Update(role)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) RoleDelete(ctx *gin.Context) {

	role := new(entity.Role)

	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.roleService.Delete(role.Rkey)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}
