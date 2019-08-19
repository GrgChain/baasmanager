package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
)

func (a *ApiController) ChaincodeAdd(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.AddChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeDeploy(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.DeployChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeUpgrade(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.UpgradeChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeQuery(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.QueryChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeInvoke(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.InvokeChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeGet(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, chain := a.chaincodeService.GetByChaincode(cc)
	if isSuccess {
		gintool.ResultOk(ctx, chain)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChaincodeList(ctx *gin.Context) {

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
	name := ctx.Query("chaincodeName")
	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, "channelId error")
		return
	}
	b, list, total := a.chaincodeService.GetList(&entity.Chaincode{
		ChaincodeName: name,
		ChannelId:     channelId,
	}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChaincodeUpdate(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chaincodeService.Update(cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeDeleted(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chaincodeService.Delete(cc.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}
