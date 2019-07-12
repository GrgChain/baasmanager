package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
	"github.com/jonluo94/baasmanager/baas-fabricengine/generate"
	"github.com/jonluo94/baasmanager/baas-fabricengine/fautil"
	"github.com/jonluo94/baasmanager/baas-core/core/fasdk"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
	"github.com/jonluo94/baasmanager/baas-fabricengine/constant"
	"github.com/jonluo94/baasmanager/baas-core/common/fileutil"
	"github.com/jonluo94/baasmanager/baas-fabricengine/config"
	"github.com/jonluo94/baasmanager/baas-core/common/json"
	"github.com/jonluo94/baasmanager/baas-core/common/util"
)

var logger = log.GetLogger("fabricengine.service", log.INFO)

type FabricService struct {
	kube KubeService
}

func NewFabricService() FabricService {
	return FabricService{
		kube: newKubeService(),
	}
}

//定义chain
func (f FabricService) defChain(ctx *gin.Context) {
	var chain model.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)

	configBuilder := generate.NewConfigBuilder(chain, paths.ArtifactPath)
	//生成crypto-feconfig.yaml
	configBuilder.BuildCryptoFile()
	//生成configtx.yaml
	configBuilder.BuildTxFile()

	artifacts := generate.NewChannelArtifacts(chain, paths.ArtifactPath)
	//生成证书文件
	artifacts.GenerateCerts()
	//生成创世区块
	artifacts.GenerateOrdererGenesis()

	gintool.ResultMsg(ctx, "success")
}

//定义channel和构建
func (f FabricService) defChannelAndBuild(ctx *gin.Context) {
	var chain model.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	configBuilder := generate.NewConfigBuilder(chain, paths.ArtifactPath)
	//生成crypto-feconfig.yaml
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
		gintool.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := fasdk.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, fautil.GetFirstOrderer(chain))
	defer fsdk.Close()
	fsdk.Setup()
	//创建channel
	fsdk.CreateChannel(fautil.GetChannelTx(chain, paths.ArtifactPath))
	//跟新锚节点
	fsdk.UpdateChannel(fautil.GetAnchorsTxs(chain, paths.ArtifactPath))
	fsdk.JoinChannel()
	/*操作fabric end*/

	gintool.ResultMsg(ctx, "success")
}

func (f FabricService) defK8sYamlAndDeploy(ctx *gin.Context) {
	var chain model.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//生成文件
	generate.NewFabricK8s(chain, paths).Build()

	datas := new(model.K8sData)

	switch chain.Consensus {
	case constant.OrdererSolo:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, f.kube.baseFiles)
	case constant.OrdererKafka:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, append(f.kube.kafkaFiles, f.kube.baseFiles...))
	case constant.OrdererEtcdraft:
		logger.Infof("no finished")

	}

	//部署k8s
	res := f.kube.deployData(datas)
	//返回
	var ret gintool.RespData
	err := json.Unmarshal(res, &ret)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if ret.Code == 0 {
		gintool.ResultMsg(ctx, "success")
	}

}

func (f FabricService) stopChainInK8s(ctx *gin.Context) {
	var chain model.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)

	datas := new(model.K8sData)
	switch chain.Consensus {
	case constant.OrdererSolo:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, f.kube.baseFiles)
	case constant.OrdererKafka:
		datas.Data = util.Yamls2Bytes(paths.K8sConfigPath, append(f.kube.kafkaFiles, f.kube.baseFiles...))
	case constant.OrdererEtcdraft:
		logger.Infof("no finished")
	}

	//停止k8s
	res := f.kube.deleteData(datas)
	//返回
	var ret gintool.RespData
	err := json.Unmarshal(res, &ret)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if ret.Code == 0 {
		gintool.ResultMsg(ctx, "success")
	}

}

func (f FabricService) releaseChain(ctx *gin.Context) {
	var chain model.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//获取目录
	err := generate.NewProjetc().RemoveProjectDir(chain)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	gintool.ResultMsg(ctx, "success")

}

func (f FabricService) buildChaincode(ctx *gin.Context) {
	var channel model.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := fasdk.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, fautil.GetFirstOrderer(chain))
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//安装chaincode
	err = fsdk.InstallChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//实例化chaincode
	resp, err := fsdk.InstantiateChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version, channel.Policy, channel.Args)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	logger.Info(string(resp))
	/*操作fabric end*/

	gintool.ResultOk(ctx, string(resp))
}

func (f FabricService) updateChaincode(ctx *gin.Context) {
	var channel model.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := fasdk.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, fautil.GetFirstOrderer(chain))
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//安装chaincode
	err = fsdk.InstallChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//升级chaincode
	resp, err := fsdk.UpgradeChaincode(channel.ChaincodeId, channel.ChaincodePath, channel.Version, channel.Policy, channel.Args)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	logger.Info(string(resp))
	/*操作fabric end*/

	gintool.ResultOk(ctx, string(resp))
}

func (f FabricService) queryChaincode(ctx *gin.Context) {
	var channel model.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := fasdk.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, fautil.GetFirstOrderer(chain))
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//创建channel
	resp, err := fsdk.QueryChaincode(channel.ChaincodeId, string(channel.Args[0]), channel.Args[1:])
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	logger.Info(string(resp))
	/*操作fabric end*/

	gintool.ResultOk(ctx, string(resp))
}

func (f FabricService) invokeChaincode(ctx *gin.Context) {
	var channel model.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	chain := channel.GetChain()
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	//连接文件
	connectConfig, err := f.getConnectConfig(chain, paths)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	/*操作fabric start*/
	fsdk := fasdk.NewFabricClient(connectConfig, chain.ChannelName, chain.PeersOrgs, fautil.GetFirstOrderer(chain))
	defer fsdk.Close()
	err = fsdk.Setup()
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//创建channel
	resp, err := fsdk.InvokeChaincode(channel.ChaincodeId, string(channel.Args[0]), channel.Args[1:])
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	logger.Info(string(resp))
	/*操作fabric end*/

	gintool.ResultOk(ctx, string(resp))
}

func (f FabricService) getConnectConfig(chain model.FabricChain, paths generate.UserChainPath) ([]byte, error) {
	nss := fautil.GetNamesapces(chain)
	res := f.kube.getChainDomain(nss)
	var ret gintool.RespData
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
	var channel model.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	fileutil.CreatedDir(fautil.GetChaincodeLocalGithub(channel))
	ioutil.WriteFile(fautil.GetChaincodeGithubFile(channel), channel.ChaincodeBytes, os.ModePerm)

	gintool.ResultOk(ctx, fautil.GetChaincodeGithub(channel))
}

func (f FabricService) downloadChaincode(ctx *gin.Context) {
	var channel model.FabricChannel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	bys, err := ioutil.ReadFile(fautil.GetChaincodeGithubFile(channel))
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	gintool.ResultOk(ctx, bys)
}

//下载artifacts
func (f FabricService) downloadArtifacts(ctx *gin.Context) {
	var chain model.FabricChain
	if err := ctx.ShouldBindJSON(&chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	//获取目录
	paths := generate.NewProjetc().BuildProjectDir(chain)
	artifactPath := paths.ArtifactPath
	tarFile := artifactPath + ".tar"
	//打包
	err := fileutil.Tar(artifactPath, tarFile, false)
	if err != nil {
		fmt.Println(err)
	}

	ctx.File(tarFile)
}

//服务
func Server() {

	fabric := NewFabricService()

	router := gin.New()
	router.Use(gintool.Logger())
	router.Use(gin.Recovery())

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

	router.Run(":" + config.Config.GetString("BaasFabricEnginePort"))
}
