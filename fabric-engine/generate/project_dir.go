package generate

import (
	"gitee.com/jonluo/baasmanager/fabric-engine/models"
	"gitee.com/jonluo/baasmanager/fabric-engine/config"
	"gitee.com/jonluo/baasmanager/fabric-engine/constant"
	"strconv"
	"gitee.com/jonluo/baasmanager/fabric-engine/util"
	"path/filepath"
	"fmt"
)

type ProjectDir struct {
}

func (p ProjectDir) BuildProjectDir(chain models.FabricChain) config.UserBaasConfig {

	artifactPath := filepath.Join(constant.BaasArtifactsDir, chain.Account, chain.ChainName)
	k8sConfig := filepath.Join(constant.BaasK8sFabricConfigDir, chain.Account, chain.ChainName)
	dataPath := filepath.Join(constant.BaasFabricDataDir, chain.Account, chain.ChainName)
	templatePath := constant.BaasK8sFabricTemplate

	util.CreatedDir(artifactPath)
	util.CreatedDir(k8sConfig)
	util.CreatedDir(dataPath)
	//创建artifact文件夹
	util.CreatedDir(filepath.Join(artifactPath, constant.ChannelArtifactsDir))

	switch chain.Consensus {
	case constant.OrdererSolo:
		domain := "orderer0." + chain.GetHostDomain(constant.OrdererSuffix)
		util.CreatedDir(filepath.Join(dataPath, domain))
	case constant.OrdererKafka:
		for i := 0; i < chain.OrderCount; i++ {
			domain := "orderer" + strconv.Itoa(i) + "." + chain.GetHostDomain(constant.OrdererSuffix)
			util.CreatedDir(filepath.Join(dataPath, domain))
		}

		for i := 0; i < 4; i++ {
			kafka := "kafka" + strconv.Itoa(i) + "." + chain.GetHostDomain(constant.KafkaSuffix)
			util.CreatedDir(filepath.Join(dataPath, kafka))
		}
	}

	for _, o := range chain.PeersOrgs {
		for i := 0; i < chain.PeerCount; i++ {
			domain := "peer" + strconv.Itoa(i) + "." + chain.GetHostDomain(o)
			util.CreatedDir(filepath.Join(dataPath, domain))
			util.CreatedDir(filepath.Join(dataPath, "couchdb."+domain))
		}
	}

	return *config.NewUserBaasConfig(artifactPath, k8sConfig, dataPath, templatePath)
}

func (p ProjectDir) RemoveProjectDir(chain models.FabricChain) error {

	artifactPath := filepath.Join(constant.BaasArtifactsDir, chain.Account, chain.ChainName)
	k8sConfig := filepath.Join(constant.BaasK8sFabricConfigDir, chain.Account, chain.ChainName)
	dataPath := filepath.Join(constant.BaasFabricDataDir, chain.Account, chain.ChainName)

	if util.RemoveDir(artifactPath) && util.RemoveDir(k8sConfig) && util.RemoveDir(dataPath){
		return nil
	}
	return fmt.Errorf("remove project dir error")

}

func NewProjetc() ProjectDir {
	return ProjectDir{}
}
