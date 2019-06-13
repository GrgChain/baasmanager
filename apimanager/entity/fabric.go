package entity

import "strings"

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

func ParseFabircChain(chain *Chain) FabricChain {
	fc := FabricChain{}
	fc.ChainName = chain.Name
	fc.Account = chain.UserAccount
	fc.Consensus = chain.Consensus
	fc.PeersOrgs = strings.Split(chain.PeersOrgs, ",")
	fc.OrderCount = chain.OrderCount
	fc.PeerCount = chain.PeerCount
	fc.ChannelName = ""
	fc.TlsEnabled = chain.TlsEnabled
	return fc
}

func ParseFabircChainAndChannel(chain *Chain, channel *Channel) FabricChain {
	fc := FabricChain{}
	fc.ChainName = chain.Name
	fc.Account = chain.UserAccount
	fc.Consensus = chain.Consensus
	fc.PeersOrgs = strings.Split(channel.Orgs, ",")
	fc.OrderCount = chain.OrderCount
	fc.PeerCount = chain.PeerCount
	fc.ChannelName = channel.ChannelName
	fc.TlsEnabled = chain.TlsEnabled
	return fc
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

func ParseFabircChannel(chain FabricChain, cc *Chaincode) FabricChannel {
	fc := FabricChannel{}
	fc.FabricChain = chain
	fc.ChaincodeId = cc.ChaincodeName
	fc.Version = cc.Version
	fc.Policy = cc.Policy
	fc.ChaincodePath = cc.GithubPath
	fc.PeerCount = chain.PeerCount
	return fc
}
