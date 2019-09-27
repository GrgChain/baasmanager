package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"net/http"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"time"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
)

func (a *ApiController) ChainAdd(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	chain.Created = time.Now().Unix()
	isSuccess, msg := a.chainService.Add(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainGet(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if a.userService.HasAdminRole(chain.UserAccount) {
		//admin 可看所有
		chain.UserAccount = ""
	}
	isSuccess, chain := a.chainService.GetByChain(chain)
	if isSuccess {
		gintool.ResultOk(ctx, chain)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChainList(ctx *gin.Context) {

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
	userAccount := ctx.Query("userAccount")

	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}

	b, list, total := a.chainService.GetList(&entity.Chain{
		Name:        name,
		UserAccount: userAccount,
	}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChainUpdate(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.Update(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainDeleted(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.Delete(chain.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainBuild(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.BuildChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}

}

func (a *ApiController) ChainRun(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.RunChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}

}

func (a *ApiController) ChainStop(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.StopChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainRelease(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.ReleaseChain(chain)
	if isSuccess {
		a.channelService.DeleteByChainId(chain.Id)
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainDownload(ctx *gin.Context) {

	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, "chainId error")
		return
	}

	chain := new(entity.Chain)
	chain.Id = chainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	reader, contentLength, name := a.chainService.DownloadChainArtifacts(chain)
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, name),
	}

	ctx.DataFromReader(http.StatusOK, contentLength, "application/x-tar", reader, extraHeaders)

}


func (a *ApiController) ChainPodsQuery(ctx *gin.Context) {

	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, "chainId error")
		return
	}

	chain := new(entity.Chain)
	chain.Id = chainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, dat := a.chainService.QueryChainPods(chain)
	if isSuccess {
		gintool.ResultOk(ctx, dat)
	} else {
		gintool.ResultFail(ctx, "query error")
	}

}

func (a *ApiController) ChangeChainResouces(ctx *gin.Context) {

	resouces := new(model.Resources)
	if err := ctx.ShouldBindJSON(resouces); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, dat := a.chainService.ChangeChainResouces(resouces)
	if isSuccess {
		gintool.ResultOk(ctx, dat)
	} else {
		gintool.ResultFail(ctx, "change error")
	}

}
