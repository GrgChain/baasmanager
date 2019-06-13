package models

import (
	"strings"
	"gitee.com/jonluo/baasmanager/fabric-engine/constant"
	"gitee.com/jonluo/baasmanager/fabric-engine/util"
	"os"
	"path/filepath"
)

type FabricChain struct {
	ChainName   string   `json:"chainName"`
	Account     string   `json:"account"`     //用户帐号
	Consensus   string   `json:"consensus"`   //共识
	PeersOrgs   []string `json:"peersOrgs"`   //参与组织 除了orderer
	OrderCount  int      `json:"orderCount"`  //orderer节点个数
	PeerCount   int      `json:"peerCount"`   //每个组织节点个数
	ChannelName string   `json:"channelName"` //channel 名
	TlsEnabled  string   `json:"tlsEnabled"`  //是否开启tls  true or false
}

func (f FabricChain) GetHostDomain(org string) string {
	return strings.ToLower(f.Account + f.ChainName + org)
}

func (f FabricChain) GetNamesapces() string {
	namespaces := make([]string, 0)
	orderer := f.GetHostDomain(constant.OrdererSuffix)
	namespaces = append(namespaces, orderer)
	for _, o := range f.PeersOrgs {
		peer := f.GetHostDomain(o)
		namespaces = append(namespaces, peer)
	}
	return strings.Join(namespaces, ",")

}

func (f FabricChain) GetChannelTx(artifactPath string) string {
	return filepath.Join(artifactPath, constant.ChannelArtifactsDir, f.ChannelName+constant.Tx)
}

func (f FabricChain) GetAnchorsTxs(artifactPath string) []string {
	anchorsTx := make([]string, len(f.PeersOrgs))
	for i, v := range f.PeersOrgs {
		tx := filepath.Join(artifactPath, constant.ChannelArtifactsDir, f.ChannelName+util.FirstUpper(v)+constant.AnchorsTx)
		anchorsTx[i] = tx
	}
	return anchorsTx
}

func (f FabricChain) GetFirstOrderer() string {
	return "orderer0." + f.GetHostDomain(constant.OrdererSuffix)
}

type FabricChannel struct {
	FabricChain

	ChaincodeId    string   `json:"chaincodeId"`
	ChaincodePath  string   `json:"chaincodePath"`
	ChaincodeBytes []byte   `json:"chaincodeBytes"`
	Version        string   `json:"version"`
	Policy         string   `json:"policy"`
	Args           [][]byte `json:"args"`
}

func (f FabricChannel) GetChain() FabricChain {
	return f.FabricChain
}

func (f FabricChannel) GetChaincodeGithub() string {
	return filepath.Join(constant.BaasChaincodeGithub, f.Account, f.ChannelName, f.ChaincodeId, f.Version)
}

func (f FabricChannel) GetChaincodeLocalGithub() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", constant.BaasChaincodeGithub, f.Account, f.ChannelName, f.ChaincodeId, f.Version)
}
func (f FabricChannel) GetChaincodeGithubFile() string {
	return filepath.Join(f.GetChaincodeLocalGithub(), constant.BaasChaincodeFile)
}
