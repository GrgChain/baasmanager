package model

import (
	"strings"
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
