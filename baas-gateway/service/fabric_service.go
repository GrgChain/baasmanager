package service

import (
	"github.com/jonluo94/baasmanager/baas-core/core/model"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
	"github.com/jonluo94/baasmanager/baas-core/common/httputil"
	"github.com/jonluo94/baasmanager/baas-gateway/config"
)

var logger = log.GetLogger("service", log.ERROR)

type FabricService struct {
}

func (g FabricService) DefChain(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/defChain", chain)
}

func (g FabricService) DefChannel(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/defChannelAndBuild", chain)
}

func (g FabricService) DeployK8sData(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/defK8sYamlAndDeploy", chain)
}

func (g FabricService) StopChain(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/stopChain", chain)
}

func (g FabricService) ReleaseChain(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/releaseChain", chain)
}

func (g FabricService) DownloadChainArtifacts(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/downloadArtifacts", chain)
}

func (g FabricService) BuildChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/buildChaincode", channel)
}
func (g FabricService) UpdateChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/updateChaincode", channel)
}

func (g FabricService) QueryChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryChaincode", channel)
}
func (g FabricService) InvokeChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/invokeChaincode", channel)
}

func (g FabricService) UploadChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/uploadChaincode", channel)
}
func (g FabricService) DownloadChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/downloadChaincode", channel)
}

func (g FabricService) QueryChainPods(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryChainPods", chain)
}

func (g FabricService) QueryLedger(channel model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryLedger", channel)
}

func (g FabricService) QueryLatestBlocks(channel model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryLatestBlocks", channel)
}

func (g FabricService) QueryBlock(channel model.FabricChain,search string) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryBlock?search="+search, channel)
}

func (g FabricService) ChangeChainPodResources(resource model.Resources) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/changeChainPodResources", resource)
}

func NewFabricService() *FabricService {
	return &FabricService{}
}
