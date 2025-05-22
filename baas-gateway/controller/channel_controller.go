package controller

import (
	"data/baas-gateway/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"strconv"
	"strings"
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

// 工具函数

func (a *ApiController) GetChannelList(orgId int) []*entity.Channel {
	org := &entity.Org{
		Id: orgId,
	}

	has, org := a.orgService.GetByOrg(org)
	if !has {
		fmt.Println("orgId not exist")
		return nil
	}

	channelNameList := strings.Split(org.Channels, ",")
	channelList := make([]*entity.Channel, len(channelNameList))
	for i := 0; i < len(channelNameList); i++ {
		channel := new(entity.Channel)
		channel.ChannelName = channelNameList[i]
		hasCh, ch := a.channelService.GetByChannel(channel)
		if !hasCh {
			fmt.Printf("channel name = %s 不存在", ch.ChannelName)
			continue
		}
		channelList[i] = ch
		fmt.Printf("ch id = %d\n ChainId = %d\n Orgs = %s\n ChannelName = %s\n UserAccount = %s\n Created = %d\n \n",
			ch.Id, ch.ChainId, ch.Orgs, ch.ChannelName, ch.UserAccount, ch.Created)
	}
	return channelList
}
