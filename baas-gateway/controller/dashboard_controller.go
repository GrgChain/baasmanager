package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
)

func (a *ApiController) DashboardCounts(ctx *gin.Context) {

	userAccount := ctx.Query("userAccount")
	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}
	isSuccess, ash := a.dashboardService.Counts(userAccount)
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) DashboardConsensusTotal(ctx *gin.Context) {

	userAccount := ctx.Query("userAccount")
	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}
	isSuccess, ash := a.dashboardService.ConsensusTotal(userAccount)
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) DashboardSevenDays(ctx *gin.Context) {

	start, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		gintool.ResultFail(ctx, "start error")
		return
	}
	end, err := strconv.Atoi(ctx.Query("end"))
	if err != nil {
		gintool.ResultFail(ctx, "end error")
		return
	}

	userAccount := ctx.Query("userAccount")
	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}
	isSuccess, ash := a.dashboardService.SevenDays(userAccount, start, end)
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}
