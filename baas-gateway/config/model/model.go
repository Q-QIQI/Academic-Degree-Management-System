package model

import "data/baas-gateway/entity"

type Dashboard struct {
	Users      int64 `json:"users"`
	Chains     int64 `json:"chains"`
	Channels   int64 `json:"channels"`
	Chaincodes int64 `json:"chaincodes"`
}

type LoginForm struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Org      int    `form:"org"      binding:"required"`
}

type OrgChannel struct {
	Channel    *entity.Channel     `json:"channel"`
	Chaincodes []*entity.Chaincode `json:"chaincodes"`
}

type InvokeArgs struct {
	ChainName     string `json:"chainName"`
	OrgName       string `json:"orgName"`
	ChannelName   string `json:"channelName"`
	ChaincodeName string `json:"chaincodeName"`
	Args          string `json:"args"`
	Fcn           string `json:"fcn"`
	Fcntype       string `json:"fcntype"`
}
