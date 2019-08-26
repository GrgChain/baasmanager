package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/xorm"
	"github.com/jonluo94/baasmanager/baas-gateway/config"
	"github.com/jonluo94/baasmanager/baas-gateway/service"
	"github.com/jonluo94/baasmanager/baas-gateway/controller"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
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
	)

	router := gin.New()
	router.Use(gintool.Logger())
	router.Use(gin.Recovery())

	gintool.UseSession(router)

	api := router.Group("/api")
	{

		api.POST("/user/login", apiController.UserLogin)
		api.POST("/user/logout", apiController.UserLogout)
		//认证校验
		api.Use(apiController.UserAuthorize)
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
		api.POST("/chaincode/invoke", apiController.ChaincodeInvoke)
		api.POST("/chaincode/get", apiController.ChaincodeGet)
		api.POST("/chaincode/delete", apiController.ChaincodeDeleted)

		api.POST("/upload", apiController.Upload)

		api.GET("/dashboard/counts", apiController.DashboardCounts)
		api.GET("/dashboard/sevenDays", apiController.DashboardSevenDays)
	}

	router.Run(":" + config.Config.GetString("BaasGatewayPort"))
}
