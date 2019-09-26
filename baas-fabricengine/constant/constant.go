package constant

const (
	//每条chain的fabric证书文件目录
	BaasArtifacts = "baas-artifacts"
	//每条chain的fabric的k8s文件目录
	BaasK8sConfig = "baas-k8s-config"
	//每条chain的fabric的数据目录
	BaasFabricData = "baas-fabric-data"
	//系统channel名称
	BaasFirstChannel = "youcannotseeme"
	//chaincode保存的文件名
	BaasChaincodeFile = "main.go"
)

const (
	//cryptogen工具生成证书目录
	CryptoConfigDir = "crypto-config"
	//configtxgen工具生成创世区块,channel交易保存的目录
	ChannelArtifactsDir = "channel-artifacts"
	//cryptogen配置
	CryptoConfigYaml = CryptoConfigDir + ".yaml"
	//configtxgen配置
	ConfigtxYaml = "configtx.yaml"
)

const (
	//configtx下的proflie
	ProfilesGenesis = "OrdererGenesis"
	ProfilesChannel = "OrgsChannel"
	//创世区块名
	GenesisBlock = "genesis.block"
	//交易后缀
	Tx = ".tx"
	//锚节点后缀
	AnchorsTx = "MSPanchors.tx"
)

//configtx文件下的配置
const (
	OrdererSuffix      = "orderer"
	OrdererMsp         = "OrdererMSP"
	OrdererSolo        = "solo"
	OrdererKafka       = "kafka"
	OrdererEtcdraft    = "etcdraft"
	KafkaSuffix        = "kafka"
	TypeImplicitMeta   = "ImplicitMeta"
	TypeSignature      = "Signature"
	RuleAnyReaders     = "ANY Readers"
	RuleAnyWriters     = "ANY Writers"
	RuleMajorityAdmins = "MAJORITY Admins"
	Country            = "CN"
	Province           = "GuangDong"
	Locality           = "GuangZhou"
)

//k8s模板文件
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
	//k8s模板的标签
	Tag = "{{%s}}"
	//私钥后缀
	PriKeySuf = "_sk"
	//msp后缀
	MspSuf = "MSP"
)
