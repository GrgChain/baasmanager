package service

import (
	"github.com/gin-gonic/gin"
	"gitee.com/jonluo/baasmanager/fabric-engine/generate"
	"net/http"
	"gitee.com/jonluo/baasmanager/fabric-engine/models"
	"gitee.com/jonluo/baasmanager/fabric-engine/constant"
	"gitee.com/jonluo/baasmanager/fabric-engine/util"
	"encoding/json"
	"github.com/jonluo94/fabric-goclient"
	"log"
	"gitee.com/jonluo/baasmanager/fabric-engine/config"
	"fmt"
	"io/ioutil"
	"os"
)

type FabricService struct {
	k8s K8sService
}

func NewFabricService() FabricService {
	return FabricService{
		k8s: K8sService{},
	}
}

//定义chain
func (f FabricService) defChain(ctx *gin.Context) {
	var chain models.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)

	configBuilder := generate.NewConfigBuilder(chain, paths.ArtifactPath)
	//生成crypto-config.yaml
	configBuilder.BuildCryptoFile()
	//生成configtx.yaml
	configBuilder.BuildTxFile()

	artifacts := generate.NewChannelArtifacts(chain, paths.ArtifactPath)
	//生成证书文件
	artifacts.GenerateCerts()
	//生成创世区块
	artifacts.GenerateOrdererGenesis()

	models.ResultMsg(ctx, "success")
}

func (f FabricService) defChannelAndBuild(ctx *gin.Context) {
	var chain models.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	configBuilder := generate.NewConfigBuilder(chain, paths.ArtifactPath)
	//生成crypto-config.yaml
	configBuilder.BuildCryptoFile()
	//生成configtx.yaml
	configBuilder.BuildTxFile()

	artifacts := generate.NewChannelArtifacts(chain, paths.ArtifactPath)
	//生成channel.tx
	artifacts.GenerateOrgsChannel()
	//生成锚节点.tx
	artifacts.GenerateAnchorPeers()

	//生成链接文件
	generate.NewConnectConfig(chain, paths.ArtifactPath).Build()
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := client.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, chain.GetFirstOrderer())
	defer fsdk.Close()
	fsdk.Setup()
	//创建channel
	fsdk.CreateChannel(chain.GetChannelTx(paths.ArtifactPath))
	//跟新锚节点
	fsdk.UpdateChannel(chain.GetAnchorsTxs(paths.ArtifactPath))
	fsdk.JoinChannel()
	/*操作fabric end*/

	models.ResultMsg(ctx, "success")
}

func (f FabricService) defK8sYamlAndDeploy(ctx *gin.Context) {
	var chain models.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//生成文件
	generate.NewFabricK8s(chain, paths).Build()

	datas := new(models.K8sData)
	switch chain.Consensus {
	case constant.OrdererSolo:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, constant.BaseFiles)
	case constant.OrdererKafka:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, append(constant.KafkaFiles, constant.BaseFiles...))
	}

	//部署k8s
	res := f.k8s.deployData(datas)

	var ret models.RespData
	err := json.Unmarshal(res, &ret)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}

	if ret.Code == 0 {
		models.ResultMsg(ctx, "success")
	}

}

func (f FabricService) stopChainInK8s(ctx *gin.Context) {
	var chain models.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)

	datas := new(models.K8sData)
	switch chain.Consensus {
	case constant.OrdererSolo:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, constant.BaseFiles)
	case constant.OrdererKafka:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, append(constant.KafkaFiles, constant.BaseFiles...))
	}

	//停止k8s
	res := f.k8s.deleteData(datas)

	var ret models.RespData
	err := json.Unmarshal(res, &ret)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}

	if ret.Code == 0 {
		models.ResultMsg(ctx, "success")
	}

}

func (f FabricService) releaseChain(ctx *gin.Context) {
	var chain models.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//获取目录
	err := generate.NewProjetc().RemoveProjectDir(chain)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}

	models.ResultMsg(ctx, "success")

}


func (f FabricService) buildChaincode(ctx *gin.Context) {
	var channel models.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := client.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, chain.GetFirstOrderer())
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//安装chaincode
	err = fsdk.InstallChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//实例化chaincode
	resp, err := fsdk.InstantiateChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version, channel.Policy, channel.Args)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	log.Println(string(resp))
	/*操作fabric end*/

	models.ResultOk(ctx, string(resp))
}

func (f FabricService) updateChaincode(ctx *gin.Context) {
	var channel models.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := client.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, chain.GetFirstOrderer())
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//安装chaincode
	err = fsdk.InstallChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//升级chaincode
	resp, err := fsdk.UpgradeChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version, channel.Policy, channel.Args)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	log.Println(string(resp))
	/*操作fabric end*/

	models.ResultOk(ctx, string(resp))
}

func (f FabricService) queryChaincode(ctx *gin.Context) {
	var channel models.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := client.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, chain.GetFirstOrderer())
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//创建channel
	resp, err := fsdk.QueryChaincode(channel.ChaincodeId, string(channel.Args[0]), channel.Args[1:])
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	log.Println(string(resp))
	/*操作fabric end*/

	models.ResultOk(ctx, string(resp))
}

func (f FabricService) invokeChaincode(ctx *gin.Context) {
	var channel models.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := client.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, chain.GetFirstOrderer())
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//创建channel
	resp, err := fsdk.InvokeChaincode(channel.ChaincodeId, string(channel.Args[0]), channel.Args[1:])
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	log.Println(string(resp))
	/*操作fabric end*/

	models.ResultOk(ctx, string(resp))
}

func (f FabricService) getConnectConfig(chain models.FabricChain, paths config.UserBaasConfig) ([]byte, error) {
	nss := chain.GetNamesapces()
	res := f.k8s.getChainDomain(nss)
	var ret models.RespData
	err := json.Unmarshal(res, &ret)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	if ret.Code != 0 {
		return nil, fmt.Errorf("%s", "no chain domans")
	}
	//连接文件
	connectConfig := generate.NewConnectConfig(chain, paths.ArtifactPath).GetBytes(ret.Data.(map[string]interface{}))
	return connectConfig, nil
}

func (f FabricService) uploadChaincode(ctx *gin.Context) {
	var channel models.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		models.ResultFail(ctx, err)
		return
	}

	util.CreatedDir(channel.GetChaincodeLocalGithub())
	ioutil.WriteFile(channel.GetChaincodeGithubFile(), channel.ChaincodeBytes, os.ModePerm)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": channel.GetChaincodeGithub()})
}

func (f FabricService) downloadChaincode(ctx *gin.Context) {
	var channel models.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		models.ResultFail(ctx, err)
		return
	}

	bys, err := ioutil.ReadFile(channel.GetChaincodeGithubFile())
	if err != nil {
		models.ResultFail(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": bys})
}

//下载artifacts
func (f FabricService) downloadArtifacts(ctx *gin.Context) {
	var chain models.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	artifactPath := paths.ArtifactPath
	tarFile := artifactPath + ".tar"
	//打包
	err := util.Tar(artifactPath, tarFile, false)
	if err != nil {
		fmt.Println(err)
	}

	ctx.File(tarFile)
}

func Server() {

	fabric := NewFabricService()

	router := gin.Default()
	router.POST("/defChain", fabric.defChain)
	router.POST("/defChannelAndBuild", fabric.defChannelAndBuild)
	router.POST("/defK8sYamlAndDeploy", fabric.defK8sYamlAndDeploy)
	router.POST("/stopChain", fabric.stopChainInK8s)
	router.POST("/releaseChain", fabric.releaseChain)
	router.POST("/buildChaincode", fabric.buildChaincode)
	router.POST("/updateChaincode", fabric.updateChaincode)
	router.POST("/queryChaincode", fabric.queryChaincode)
	router.POST("/invokeChaincode", fabric.invokeChaincode)
	router.POST("/uploadChaincode", fabric.uploadChaincode)
	router.POST("/downloadChaincode", fabric.downloadChaincode)
	router.POST("/downloadArtifacts", fabric.downloadArtifacts)
	router.Run(":4991")
}
