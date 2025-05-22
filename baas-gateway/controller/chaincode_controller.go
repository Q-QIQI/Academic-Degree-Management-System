package controller

import (
	"data/baas-gateway/entity"
	"data/baas-gateway/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"strconv"
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

func (a *ApiController) ChaincodeLedgerQuery(ctx *gin.Context) {

	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, "channelId error")
		return
	}

	channel := new(entity.Channel)
	channel.Id = channelId
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

	isSuccess, msg := a.chaincodeService.QueryLedger(chain, channel)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeLatestBlocksQuery(ctx *gin.Context) {

	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, "channelId error")
		return
	}

	channel := new(entity.Channel)
	channel.Id = channelId
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

	isSuccess, msg := a.chaincodeService.QueryLatestBlocks(chain, channel)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeBlockQuery(ctx *gin.Context) {

	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, "channelId error")
		return
	}
	search := ctx.Query("search")

	channel := new(entity.Channel)
	channel.Id = channelId
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

	isSuccess, msg := a.chaincodeService.QueryBlock(chain, channel, search)
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

func (a *ApiController) QueryChaincodeByOrgId(ctx *gin.Context) {
	invokeArgs := new(model.InvokeArgs)
	if err := ctx.ShouldBindJSON(invokeArgs); err != nil {
		fmt.Println("QueryChaincodeByOrgId bind json error", err.Error())
		gintool.ResultFail(ctx, err)
		return
	}
	fmt.Printf("ChainName <%s> OrgName <%s> ChannelName <%s> ChaincodeName <%s>\n", invokeArgs.ChainName, invokeArgs.OrgName, invokeArgs.ChannelName, invokeArgs.ChaincodeName)
	fmt.Printf("Fcn %s\n", invokeArgs.Fcn)
	fmt.Printf("Fcntype %s\n", invokeArgs.Fcntype)
	fmt.Printf("Args %s\n", invokeArgs.Args)
	chaincodes := a.chaincodeService.GetChaincodeByNames(
		invokeArgs.ChainName,
		invokeArgs.OrgName,
		invokeArgs.ChannelName,
		invokeArgs.ChaincodeName,
	)
	fmt.Printf(">> debug Chaincode Id <%d>, ChaincodeName <%s> ChannelId <%d>\n", chaincodes[0].Id, chaincodes[0].ChaincodeName, chaincodes[0].ChannelId)
	result := a.QueryChaincodeByInvokeArgs(
		chaincodes[0],
		invokeArgs.Args,
		invokeArgs.Fcn,
		invokeArgs.Fcntype,
	)
	gintool.ResultOk(ctx, result)
}

func (a *ApiController) InvokeChaincodeByOrgId(ctx *gin.Context) {
	invokeArgs := new(model.InvokeArgs)
	if err := ctx.ShouldBindJSON(invokeArgs); err != nil {
		fmt.Println("InvokeChaincodeByOrgId bind json error")
		gintool.ResultFail(ctx, err)
		return
	}
	fmt.Printf("ChainName <%s> OrgName <%s> ChannelName <%s> ChaincodeName <%s>\n", invokeArgs.ChainName, invokeArgs.OrgName, invokeArgs.ChannelName, invokeArgs.ChaincodeName)
	fmt.Printf("Fcn %s\n", invokeArgs.Fcn)
	fmt.Printf("Fcntype %s\n", invokeArgs.Fcntype)
	fmt.Printf("Args %s\n", invokeArgs.Args)
	chaincodes := a.chaincodeService.GetChaincodeByNames(
		invokeArgs.ChainName,
		invokeArgs.OrgName,
		invokeArgs.ChannelName,
		invokeArgs.ChaincodeName,
	)
	result := a.InvokeChaincodeByInvokeArgs(
		chaincodes[0],
		invokeArgs.Args,
		invokeArgs.Fcn,
		invokeArgs.Fcntype,
	)
	gintool.ResultOk(ctx, result)
}

// 工具函数

func (a *ApiController) InvokeChaincodeByInvokeArgs(chaincode entity.Chaincode, args string, fcn string, fcntype string) string {
	cc := chaincode
	cc.Args = args
	channel := &entity.Channel{
		Id: cc.ChannelId,
	}
	fmt.Printf("before GetByChannelId channel Id = <%d>\n", channel.Id)
	_, channel = a.channelService.GetByChannelId(channel)

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	_, chain = a.chainService.GetByChain(chain)
	fmt.Printf("InvokeChaincodeByInvokeArgs chain id = %d name = %s\n", chain.Id, chain.Name)

	_, msg := a.chaincodeService.InvokeChaincode(chain, channel, &cc)
	return msg
}

func (a *ApiController) QueryChaincodeByInvokeArgs(chaincode entity.Chaincode, args string, fcn string, fcntype string) string {
	cc := chaincode
	cc.Args = args
	//channel := new(entity.Channel)
	//channel.Id = cc.ChannelId
	channel := &entity.Channel{
		Id: cc.ChannelId,
	}
	fmt.Printf("before GetByChannelId channel Id = <%d>\n", channel.Id)
	_, channel = a.channelService.GetByChannelId(channel)

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	_, chain = a.chainService.GetByChain(chain)
	fmt.Printf("QueryChaincodeByInvokeArgs chain id = %d name = %s\n", chain.Id, chain.Name)

	_, msg := a.chaincodeService.QueryChaincode(chain, channel, &cc)
	return msg
}
