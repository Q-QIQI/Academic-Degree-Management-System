package controller

import (
	"fmt"
	"strings"

	"data/baas-gateway/entity"
	"data/baas-gateway/model"
	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-core/common/password"
)

//type OrgLoginResult struct {
//	token *entity.JwtToken
//	channels *[]model.OrgChannel
//}

func (a *ApiController) OrgLogin(ctx *gin.Context) {
	login := new(model.LoginForm)
	if err := ctx.ShouldBind(&login); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	user := &entity.User{
		Account: login.UserName,
	}
	has, u := a.userService.GetByUser(user)
	if !has {
		gintool.ResultFail(ctx, "username error")
		return
	}
	vali := password.Validate(login.Password, u.Password)
	if !vali {
		gintool.ResultFail(ctx, "password error")
		return
	}

	type UserInfo map[string]interface{}
	// orgid to channel, chain, chaincode
	org := &entity.Org{
		Id: login.Org,
	}
	has, org = a.orgService.GetByOrg(org)
	if has {
		fmt.Println("has", has)
		fmt.Printf("org.Id=%d\n org.Chain=%d\n org.Channels=%s\n org.Name=%s\n org.Publickey=%s\n", org.Id, org.Chain, org.Channels, org.Name, org.Publickey)
	}

	token := a.userService.GetToken(u)
	//保存session
	gintool.SetSession(ctx, token.Token, u.Id)
	gintool.ResultOk(ctx, token)
}

func (a *ApiController) OrgAdd(ctx *gin.Context) {
	org := new(entity.Org)

	if err := ctx.ShouldBindJSON(org); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.orgService.Add(org)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) OrgDelete(ctx *gin.Context) {
	org := new(entity.Org)

	if err := ctx.ShouldBindJSON(org); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.orgService.Delete(org)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) OrgUpdateByChainAndOrgName(ctx *gin.Context) {
	fmt.Println(">>> OrgUpdateByChainAndOrgName")
	org := new(entity.Org)

	if err := ctx.ShouldBindJSON(org); err != nil {
		fmt.Printf("--- OrgUpdateByChainAndOrgName bind json error")
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.orgService.Update(org)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) OrgUpdateById(ctx *gin.Context) {
	org := new(entity.Org)

	if err := ctx.ShouldBindJSON(org); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.orgService.UpdateById(org)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) OrgAll(ctx *gin.Context) {
	isSuccess, data := a.orgService.GetAllList()
	if isSuccess {
		gintool.ResultOk(ctx, data)
	} else {
		gintool.ResultFail(ctx, data)
	}
}

func (a *ApiController) GetChannelListByOrgId(ctx *gin.Context) {
	org := new(entity.Org)

	if err := ctx.ShouldBindJSON(org); err != nil {
		fmt.Println("--- ctx.ShouldBindJSON(org) error")
		gintool.ResultFail(ctx, err)
		return
	}
	channelList := a.GetChannelList(org.Id)
	gintool.ResultOk(ctx, channelList)
}

func (a *ApiController) GetOrgChannelAndChaincodeListByOrgId(ctx *gin.Context) {
	org := new(entity.Org)

	if err := ctx.ShouldBindJSON(org); err != nil {
		fmt.Println("--- ctx.ShouldBindJSON(org) error")
		gintool.ResultFail(ctx, err)
		return
	}
	ocList := a.GetOrgChannelByOrgId(org.Id)
	gintool.ResultOk(ctx, ocList)
}

func (a *ApiController) GetOrgChannelByOrgId(orgId int) []*model.OrgChannel {
	org := &entity.Org{
		Id: orgId,
	}

	has, org := a.orgService.GetByOrg(org)
	if !has {
		fmt.Println("--- a.orgService.GetByOrg error")
		return nil
	}

	fmt.Printf("org.Channels: %s\n", org.Channels)
	channelNameList := strings.Split(org.Channels, ",")
	channelList := make([]*entity.Channel, len(channelNameList))
	ocList := make([]*model.OrgChannel, len(channelNameList))
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

		chaincode := new(entity.Chaincode)
		chaincode.ChannelId = ch.Id
		hasCC, list, _ := a.chaincodeService.GetList(chaincode, 1, 10)

		oc := new(model.OrgChannel)
		oc.Channel = ch
		if !hasCC {
			fmt.Printf("channel id = %s 的链码不存在", ch.Id)
			oc.Chaincodes = nil
		} else {
			oc.Chaincodes = list
		}
		ocList[i] = oc
	}
	return ocList
}
