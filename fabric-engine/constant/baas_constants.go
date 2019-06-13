package constant

import (
	"os"
	"log"
	"path/filepath"
)

var (
	BaasRootPath  string
	BaasNfsServer string
	BaasK8sEngine string

	BaasArtifactsDir       string
	BaasK8sFabricConfigDir string
	BaasFabricDataDir      string
	BaasK8sFabricTemplate  string

	BaseFiles  []string
	KafkaFiles []string
)

const (
	BaasChaincodeGithub = "github.com/baaschaincodes"
	BaasChaincodeFile   = "main.go"
	BaasFirstChannel    = "youcannotseeme"

	BaasNfsShared  = "baas-nfs-shared"
	BaasArtifacts  = "baas-artifacts"
	BaasK8sConfig  = "baas-k8s-config"
	BaasFabricData = "baas-fabric-data"
	BaasTemplate   = "baas-template"
)
const (
	CryptoConfigDir     = "crypto-config"
	ChannelArtifactsDir = "channel-artifacts"
	CryptoConfigYaml    = "crypto-config.yaml"
	ConfigtxYaml        = "configtx.yaml"
)
const (
	K8sNfsYaml       = "nfs.yaml"
	K8sNamespaceYaml = "namespace.yaml"
	K8sOrdererYaml   = "orderer.yaml"
	K8sPeerYaml      = "peer.yaml"
	K8sCaYaml        = "ca.yaml"
	K8sCliYaml       = "cli.yaml"
	K8sZookeeperYaml = "zookeeper.yaml"
	K8sKafkaYaml     = "kafka.yaml"
)
const (
	OrdererBatchTimeout      = "2s"
	OrdererMaxMessageCount   = 500
	OrdererAbsoluteMaxBytes  = "99 MB"
	OrdererPreferredMaxBytes = "512 KB"
)

func init() {
	//BaseFiles = []string{K8sNamespaceYaml, K8sNfsYaml, K8sOrdererYaml, K8sPeerYaml, K8sCaYaml, K8sCliYaml}
	BaseFiles = []string{K8sNamespaceYaml, K8sNfsYaml, K8sOrdererYaml, K8sPeerYaml, K8sCaYaml}
	KafkaFiles = []string{K8sZookeeperYaml, K8sKafkaYaml}

	BaasRootPath = os.Getenv("BaasRootPath")
	BaasNfsServer = os.Getenv("BaasNfsServer")
	BaasK8sEngine = os.Getenv("BaasK8sEngine")

	BaasArtifactsDir = filepath.Join(BaasRootPath, BaasNfsShared, BaasArtifacts)
	BaasK8sFabricConfigDir = filepath.Join(BaasRootPath, BaasNfsShared, BaasK8sConfig)
	BaasFabricDataDir = filepath.Join(BaasRootPath, BaasNfsShared, BaasFabricData)
	BaasK8sFabricTemplate = filepath.Join(BaasRootPath, BaasTemplate)

	log.Println("BaasRootPath:", BaasRootPath)
	log.Println("BaasNfsServer:", BaasNfsServer)
	log.Println("BaasK8sEngine:", BaasK8sEngine)
	if BaasRootPath == "" || BaasNfsServer == "" || BaasK8sEngine == "" {
		log.Fatal("no env BaasRootPath or BaasNfsServer or BaasK8sEngine")
	}
}
