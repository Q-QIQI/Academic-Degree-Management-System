package main

import (
	"data/baas-gateway/config"
	"data/baas-gateway/controller"

	dataCtr "data/baas-gateway/controller/data"
	"data/baas-gateway/service"
	"data/baas-gateway/service/data"

	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-core/common/xorm"
)

func main() {
	dbengine := xorm.GetEngine(config.Config.GetString("BaasGatewayDbconfig"))
	fabricService := service.NewFabricService()
	apiController := controller.NewApiController(
		service.NewUserService(dbengine),
		service.NewRoleService(dbengine),
		service.NewChainService(dbengine, fabricService),
		service.NewChannelService(dbengine, fabricService),
		service.NewChaincodeService(dbengine, fabricService),
		service.NewDashboardService(dbengine),
		service.NewOrgService(dbengine, fabricService),
	)

	router := gin.New()
	router.Use(gintool.Logger())
	router.Use(gin.Recovery())
	router.Use(apiController.Cors)

	gintool.UseSession(router)

	dataController := dataCtr.NewDataControl(data.NewDataService(dbengine, fabricService))
	dataApi := router.Group("/api")

	dataApi.GET("/trace", dataController.FindTrace)
	{
		dataApi.Use(apiController.UserAuthorize)
		dataApi.GET("/educational/infor/list", dataController.FindEducationalInforList)
		dataApi.GET("/educational/infor/list/self", dataController.FindEducationalInforListSelf)
		dataApi.GET("/educational/infor/one", dataController.GetEducationalInforByUserId)
		dataApi.POST("/educational/infor/add", dataController.CreateEducationalInfor)
		dataApi.POST("/educational/infor/update", dataController.UpdateEducationalInforById)

		//Company
		dataApi.GET("/company/infor/list", dataController.FindCompanyList)
		dataApi.GET("/company/infor/one", dataController.GetCompanyByUserId)
		dataApi.POST("/company/infor/add", dataController.CreateCompany)
		dataApi.POST("/company/infor/update", dataController.UpdateCompanyById)

		//EducationalApplication
		dataApi.GET("/educational/application/list", dataController.FindEducationalApplicationList)
		dataApi.POST("/educational/application/add", dataController.CreateEducationalApplication)
		dataApi.POST("/educational/application/update", dataController.UpdateEducationalApplicationById)

	}

	api := router.Group("/api")
	{

		// 用户信息
		api.GET("/intellectual/users", apiController.UserAuthorize, apiController.FindUserList)
		// 更新用户信息
		api.POST("/intellectual/user/update", apiController.UserAuthorize, apiController.UpdateUser)
		// 上传文件
		api.POST("/file/upload", apiController.UpFile)
		api.GET("/file/download", apiController.DownFile)
		api.POST("/user/login", apiController.UserLogin)
		api.POST("/user/logout", apiController.UserLogout)
		api.POST("/user/register", apiController.UserRegister)
		// 认证校验
		// api.Use(apiController.UserAuthorize)
		api.GET("/user/info", apiController.UserInfo)
		api.GET("/user/list", apiController.UserList)
		api.POST("/user/add", apiController.UserAdd)
		api.POST("/user/addAuth", apiController.UserAddAuth)
		api.POST("/user/delAuth", apiController.UserDelAuth)
		api.POST("/user/update", apiController.UserUpdate)
		api.POST("/user/delete", apiController.UserDelete)

		api.GET("/role/list", apiController.RoleList)
		api.GET("/role/allList", apiController.RoleAllList)
		api.POST("/role/add", apiController.RoleAdd)
		api.POST("/role/update", apiController.RoleUpdate)
		api.POST("/role/delete", apiController.RoleDelete)

		api.GET("/chain/list", apiController.ChainList)
		api.POST("/chain/add", apiController.ChainAdd)
		api.POST("/chain/update", apiController.ChainUpdate)
		api.POST("/chain/get", apiController.ChainGet)
		api.POST("/chain/delete", apiController.ChainDeleted)
		api.POST("/chain/build", apiController.ChainBuild)
		api.POST("/chain/run", apiController.ChainRun)
		api.POST("/chain/stop", apiController.ChainStop)
		api.POST("/chain/release", apiController.ChainRelease)
		api.POST("/chain/changeSize", apiController.ChangeChainResouces)
		api.GET("/chain/download", apiController.ChainDownload)
		api.GET("/chain/podsQuery", apiController.ChainPodsQuery)

		api.POST("/channel/add", apiController.ChannelAdd)
		api.POST("/channel/get", apiController.ChannelGet)
		api.GET("/channel/allList", apiController.ChannelAll)

		api.GET("/chaincode/list", apiController.ChaincodeList)
		api.POST("/chaincode/add", apiController.ChaincodeAdd)
		api.POST("/chaincode/deploy", apiController.ChaincodeDeploy)
		api.POST("/chaincode/upgrade", apiController.ChaincodeUpgrade)
		api.POST("/chaincode/query", apiController.ChaincodeQuery)
		api.GET("/chaincode/queryLedger", apiController.ChaincodeLedgerQuery)
		api.GET("/chaincode/queryLatestBlocks", apiController.ChaincodeLatestBlocksQuery)
		api.GET("/chaincode/queryBlock", apiController.ChaincodeBlockQuery)
		api.POST("/chaincode/invoke", apiController.ChaincodeInvoke)
		api.POST("/chaincode/get", apiController.ChaincodeGet)
		api.POST("/chaincode/delete", apiController.ChaincodeDeleted)

		api.POST("/upload", apiController.Upload)
		api.POST("/uploadToNano", apiController.UploadToNano)

		api.GET("/dashboard/counts", apiController.DashboardCounts)
		api.GET("/dashboard/sevenDays", apiController.DashboardSevenDays)
		api.GET("/dashboard/consensusTotal", apiController.DashboardConsensusTotal)

		api.GET("/org/list", apiController.OrgAll)
		api.POST("/org/add", apiController.OrgAdd)
		api.POST("/org/delete", apiController.OrgDelete)
		api.POST("/org/update", apiController.OrgUpdateByChainAndOrgName)

		api.POST("/org/chaincode/list", apiController.GetChannelListByOrgId)
		api.POST("/org/channel/chaincode/list", apiController.GetOrgChannelAndChaincodeListByOrgId)

		api.POST("/connector/chaincode/query", apiController.QueryChaincodeByOrgId)
		api.POST("/connector/chaincode/invoke", apiController.InvokeChaincodeByOrgId)

	}

	router.Run(":" + config.Config.GetString("BaasGatewayPort"))
}
