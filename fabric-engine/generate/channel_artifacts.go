package generate

import (
	"sync"
	"gitee.com/jonluo/baasmanager/fabric-engine/constant"
	"gitee.com/jonluo/baasmanager/fabric-engine/models"
	"gitee.com/jonluo/baasmanager/fabric-engine/tools"
	"gitee.com/jonluo/baasmanager/fabric-engine/util"
	"path/filepath"
)

type ChannelArtifacts struct {
	channelName         string
	cryptoConfigFile    string
	cryptoConfigDir     string
	channelArtifactsDir string
	rootPath            string
	anchorOrgs          []string
}

func (c *ChannelArtifacts) GenerateCerts() {
	gen := tools.NewCryptogen(c.cryptoConfigFile, c.cryptoConfigDir)
	gen.Exec("generate")
}

func (c *ChannelArtifacts) GenerateOrdererGenesis() {
	txgen := tools.NewConfigtxgen()
	txgen.SetConfigPath(c.rootPath)
	txgen.SetProfile(constant.ProfilesGenesis)
	txgen.SetChannelID(constant.BaasFirstChannel)
	txgen.SetOutputBlock(filepath.Join(c.channelArtifactsDir, constant.GenesisBlock))
	txgen.Exec()
}

func (c *ChannelArtifacts) GenerateOrgsChannel() {
	txgen := tools.NewConfigtxgen()
	txgen.SetConfigPath(c.rootPath)
	txgen.SetProfile(constant.ProfilesChannel)
	txgen.SetOutputChannelCreateTx(filepath.Join(c.channelArtifactsDir, c.channelName+constant.Tx))
	txgen.SetChannelID(c.channelName)
	txgen.Exec()
}
func (c *ChannelArtifacts) GenerateAnchorPeers() {
	var wg sync.WaitGroup

	for _, v := range c.anchorOrgs {
		wg.Add(1)
		go func(o string) {

			org := util.FirstUpper(o)
			txgen := tools.NewConfigtxgen()
			txgen.SetConfigPath(c.rootPath)
			txgen.SetProfile(constant.ProfilesChannel)
			txgen.SetOutputAnchorPeersUpdate(filepath.Join(c.channelArtifactsDir, c.channelName+org+constant.AnchorsTx))
			txgen.SetChannelID(c.channelName)
			txgen.SetAsOrg(org)
			txgen.Exec()

			wg.Done()
		}(v)

	}
	wg.Wait()
}

func NewChannelArtifacts(chain models.FabricChain, rootPath string) *ChannelArtifacts {
	return &ChannelArtifacts{
		channelName:         chain.ChannelName,
		cryptoConfigFile:    filepath.Join(rootPath, constant.CryptoConfigYaml),
		cryptoConfigDir:     filepath.Join(rootPath, constant.CryptoConfigDir),
		channelArtifactsDir: filepath.Join(rootPath, constant.ChannelArtifactsDir),
		rootPath:            rootPath,
		anchorOrgs:          chain.PeersOrgs,
	}
}
