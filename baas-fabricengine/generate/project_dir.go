package generate

import (
	"strconv"
	"path/filepath"
	"fmt"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
	"github.com/jonluo94/baasmanager/baas-fabricengine/constant"
	"github.com/jonluo94/baasmanager/baas-core/common/fileutil"
	"github.com/jonluo94/baasmanager/baas-fabricengine/config"
)

//用户chain路径
type UserChainPath struct {
	ArtifactPath  string
	K8sConfigPath string
	DataPath      string
	TemplatePath  string
}

func NewUserChainPath(artifactPath, k8sConfig, dataPath, templatePath string) UserChainPath {
	return UserChainPath{
		ArtifactPath:  artifactPath,
		K8sConfigPath: k8sConfig,
		DataPath:      dataPath,
		TemplatePath:  templatePath,
	}
}

//工程目录
type ProjectDir struct {
	BaasArtifactsDir       string
	BaasK8sFabricConfigDir string
	BaasFabricDataDir      string
}

func (p ProjectDir) BuildProjectDir(chain model.FabricChain) UserChainPath {
	//nfs shared
	artifactPath := filepath.Join(p.BaasArtifactsDir, chain.Account, chain.ChainName)
	k8sConfig := filepath.Join(p.BaasK8sFabricConfigDir, chain.Account, chain.ChainName)
	dataPath := filepath.Join(p.BaasFabricDataDir, chain.Account, chain.ChainName)
	//模板
	templatePath := filepath.Join(config.Config.GetString("BaasRootPath"), config.Config.GetString("BaasTemplate"))

	fileutil.CreatedDir(artifactPath)
	fileutil.CreatedDir(k8sConfig)
	fileutil.CreatedDir(dataPath)
	//创建artifact文件夹
	fileutil.CreatedDir(filepath.Join(artifactPath, constant.ChannelArtifactsDir))

	switch chain.Consensus {
	case constant.OrdererSolo:
		domain := "orderer0." + chain.GetHostDomain(constant.OrdererSuffix)
		fileutil.CreatedDir(filepath.Join(dataPath, domain))
	case constant.OrdererKafka:
		for i := 0; i < chain.OrderCount; i++ {
			domain := "orderer" + strconv.Itoa(i) + "." + chain.GetHostDomain(constant.OrdererSuffix)
			fileutil.CreatedDir(filepath.Join(dataPath, domain))
		}

		for i := 0; i < 4; i++ {
			kafka := "kafka" + strconv.Itoa(i) + "." + chain.GetHostDomain(constant.KafkaSuffix)
			fileutil.CreatedDir(filepath.Join(dataPath, kafka))
		}
	case constant.OrdererEtcdraft:
		for i := 0; i < chain.OrderCount; i++ {
			domain := "orderer" + strconv.Itoa(i) + "." + chain.GetHostDomain(constant.OrdererSuffix)
			fileutil.CreatedDir(filepath.Join(dataPath, domain))
		}
	}

	for _, o := range chain.PeersOrgs {
		for i := 0; i < chain.PeerCount; i++ {
			domain := "peer" + strconv.Itoa(i) + "." + chain.GetHostDomain(o)
			fileutil.CreatedDir(filepath.Join(dataPath, domain))
			fileutil.CreatedDir(filepath.Join(dataPath, "couchdb."+domain))
		}
	}

	return NewUserChainPath(artifactPath, k8sConfig, dataPath, templatePath)
}

func (p ProjectDir) RemoveProjectDir(chain model.FabricChain) error {

	artifactPath := filepath.Join(p.BaasArtifactsDir, chain.Account, chain.ChainName)
	k8sConfig := filepath.Join(p.BaasK8sFabricConfigDir, chain.Account, chain.ChainName)
	dataPath := filepath.Join(p.BaasFabricDataDir, chain.Account, chain.ChainName)

	if fileutil.RemoveDir(artifactPath) && fileutil.RemoveDir(k8sConfig) && fileutil.RemoveDir(dataPath) {
		return nil
	}
	return fmt.Errorf("remove project dir error")

}

func NewProjetc() ProjectDir {
	baasNfsSharedDir := filepath.Join(config.Config.GetString("BaasRootPath"), config.Config.GetString("BaasNfsShared"))
	baasArtifactsDir := filepath.Join(baasNfsSharedDir, constant.BaasArtifacts)
	BaasK8sFabricConfigDir := filepath.Join(baasNfsSharedDir, constant.BaasK8sConfig)
	BaasFabricDataDir := filepath.Join(baasNfsSharedDir, constant.BaasFabricData)
	return ProjectDir{
		BaasArtifactsDir:       baasArtifactsDir,
		BaasK8sFabricConfigDir: BaasK8sFabricConfigDir,
		BaasFabricDataDir:      BaasFabricDataDir,
	}
}
