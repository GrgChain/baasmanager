package service

import (
	"gitee.com/jonluo/baasmanager/apimanager/config"
	"gitee.com/jonluo/baasmanager/apimanager/entity"
	"github.com/jonluo94/commontools/httputil"
	"github.com/jonluo94/commontools/log"
)

var logger = log.GetLogger("service", log.ERROR)

type FabricService struct {
}

func (g FabricService) DefChain(chain entity.FabricChain) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/defChain", chain)
}

func (g FabricService) DefChannel(chain entity.FabricChain) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/defChannelAndBuild", chain)
}

func (g FabricService) DeployK8sData(chain entity.FabricChain) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/defK8sYamlAndDeploy", chain)
}

func (g FabricService) StopChain(chain entity.FabricChain) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/stopChain", chain)
}

func (g FabricService) ReleaseChain(chain entity.FabricChain) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/releaseChain", chain)
}

func (g FabricService) DownloadChainArtifacts(chain entity.FabricChain) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/downloadArtifacts", chain)
}

func (g FabricService) BuildChaincode(channel entity.FabricChannel) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/buildChaincode", channel)
}
func (g FabricService) UpdateChaincode(channel entity.FabricChannel) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/updateChaincode", channel)
}

func (g FabricService) QueryChaincode(channel entity.FabricChannel) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/queryChaincode", channel)
}
func (g FabricService) InvokeChaincode(channel entity.FabricChannel) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/invokeChaincode", channel)
}

func (g FabricService) UploadChaincode(channel entity.FabricChannel) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/uploadChaincode", channel)
}
func (g FabricService) DownloadChaincode(channel entity.FabricChannel) []byte {
	return httputil.PostJson(config.BaasFabricEngine+"/downloadChaincode", channel)
}

func NewFabricService() *FabricService {
	return &FabricService{}
}
