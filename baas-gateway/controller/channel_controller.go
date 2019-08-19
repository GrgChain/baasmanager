package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
)

func (a *ApiController) ChannelAdd(ctx *gin.Context) {

	channel := new(entity.Channel)

	if err := ctx.ShouldBindJSON(channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.channelService.AddChannel(chain, channel)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChannelGet(ctx *gin.Context) {

	chn := new(entity.Channel)

	if err := ctx.ShouldBindJSON(chn); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, chn := a.channelService.GetByChannel(chn)
	if isSuccess {
		gintool.ResultOk(ctx, chn)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChannelAll(ctx *gin.Context) {

	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, "chainId error")
		return
	}
	isSuccess, data := a.channelService.GetAllList(chainId)
	if isSuccess {
		gintool.ResultOk(ctx, data)
	} else {
		gintool.ResultFail(ctx, data)
	}
}
