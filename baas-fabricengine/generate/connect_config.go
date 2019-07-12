package generate

import (
	"bytes"
	"strings"
	"fmt"
	"strconv"
	"io/ioutil"
	"os"
	"encoding/json"
	"path/filepath"
	"github.com/jonluo94/baasmanager/baas-fabricengine/constant"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
	"github.com/jonluo94/baasmanager/baas-core/common/util"
	"github.com/jonluo94/baasmanager/baas-core/core/fasdk"
)

var line = `
`

var (
	org          = fmt.Sprintf(constant.Tag, "org")
	msp          = fmt.Sprintf(constant.Tag, "msp")
	domain       = fmt.Sprintf(constant.Tag, "domain")
	node         = fmt.Sprintf(constant.Tag, "node")
	nodeType     = fmt.Sprintf(constant.Tag, "type")
	noTls        = fmt.Sprintf(constant.Tag, "no-tls")
	cryptoConfig = fmt.Sprintf(constant.Tag, "crypto-kubeconfig")
	channel      = fmt.Sprintf(constant.Tag, "channel")
	user         = fmt.Sprintf(constant.Tag, "user")
)
//fasdk链接文件
type ConnectConfig struct {
	model.FabricChain
	cryptoConfig string //证书
	configYaml   string //配置文件
	buffer       *bytes.Buffer
}

//client节点
func (c *ConnectConfig) setClient() {
	client := `
client:
  organization: {{org}}
  logging:
    level: info
  cryptoconfig:
    path: {{crypto-kubeconfig}}
  credentialStore:
    path: /tmp/state-store
    cryptoStore:
      path: /tmp/msp
  tlsCerts:
    systemCertPool: true
`
	client = strings.Replace(client, org, c.PeersOrgs[0], -1)
	client = strings.Replace(client, cryptoConfig, c.cryptoConfig, -1)
	c.buffer.WriteString(client)

}

//channel
func (c *ConnectConfig) setChannels() {
	channels := `
channels:
  {{channel}}:
    peers:
`
	peers := ``
	for _, org := range c.PeersOrgs {
		for i := 0; i < c.PeerCount; i++ {

			peer := `      peer` + strconv.Itoa(i) + `.` + c.GetHostDomain(org) + `:` + line

			peers += peer
		}
	}
	channels += peers
	channels = strings.Replace(channels, channel, c.ChannelName, -1)
	c.buffer.WriteString(channels)

}

//organizations
func (c *ConnectConfig) setOrganizations() {
	organizations := `
organizations:
`
	for _, o := range c.PeersOrgs {
		organ := `
  {{org}}:
    mspid: {{msp}}
    cryptoPath:  peerOrganizations/{{domain}}/users/{username}@{{domain}}/msp
    peers:
`
		organ = strings.Replace(organ, org, o, -1)
		organ = strings.Replace(organ, msp, util.FirstUpper(o)+constant.MspSuf, -1)

		host := c.GetHostDomain(o)

		peers := ``
		for i := 0; i < c.PeerCount; i++ {
			peer := `      - peer` + strconv.Itoa(i) + `.` + host + line
			peers += peer
		}

		organ += peers
		organ += `
    certificateAuthorities:
      - ca.{{domain}}
`
		organ = strings.Replace(organ, domain, host, -1)
		organizations += organ
	}
	c.buffer.WriteString(organizations)

}

//orderer
func (c *ConnectConfig) setOrderers() {
	orderers := `
orderers:
`
	host := c.GetHostDomain(constant.OrdererSuffix)
	for i := 0; i < c.OrderCount; i++ {
		orderNode := "orderer" + strconv.Itoa(i) + "." + host
		orderer := setNode("orderer", host, orderNode, c.cryptoConfig, c.TlsEnabled)
		orderers += orderer
	}

	c.buffer.WriteString(orderers)
}

//peer
func (c *ConnectConfig) setPeers() {

	peers := `
peers:
`
	for _, o := range c.PeersOrgs {
		host := c.GetHostDomain(o)
		for i := 0; i < c.PeerCount; i++ {
			peerNode := "peer" + strconv.Itoa(i) + "." + host
			peer := setNode("peer", host, peerNode, c.cryptoConfig, c.TlsEnabled)
			peers += peer
		}

	}
	c.buffer.WriteString(peers)
}

//ca
func (c *ConnectConfig) setCertificateAuthorities() {
	cas := `
certificateAuthorities:
`
	for _, o := range c.PeersOrgs {
		host := c.GetHostDomain(o)
		ca := `
  ca.{{domain}}:
    url: {{no-tls}}://ca.{{domain}}
    tlsCACerts:
      path: {{crypto-kubeconfig}}/peerOrganizations/{{domain}}/tlsca/tlsca.{{domain}}-cert.pem
      client:
        key:
          path: {{crypto-kubeconfig}}/peerOrganizations/{{domain}}/users/{{user}}@{{domain}}/tls/client.key
        cert:
          path: {{crypto-kubeconfig}}/peerOrganizations/{{domain}}/users/{{user}}@{{domain}}/tls/client.crt
    registrar:
      enrollId: admin
      enrollSecret: adminpw
`
		ca = strings.Replace(ca, domain, host, -1)
		ca = strings.Replace(ca, user, fasdk.User, -1)
		ca = strings.Replace(ca, cryptoConfig, c.cryptoConfig, -1)

		tls := c.TlsEnabled
		if tls == "true" {
			tls = "https"
		} else {
			tls = "http"
		}

		ca = strings.Replace(ca, noTls, tls, -1)
		cas += ca
	}

	c.buffer.WriteString(cas)
}

//实体代理
func (c *ConnectConfig) setEntityMatchers() {
	entityMatchers := `
entityMatchers:
`

	entity := `
    - pattern: {{node}}
      urlSubstitutionExp: {{{{node}}}}
      sslTargetOverrideUrlSubstitutionExp: {{node}}
      mappedHost: {{node}}
`

	peers := `
  peer:`
	for _, o := range c.PeersOrgs {
		orgHost := c.GetHostDomain(o)
		for i := 0; i < c.PeerCount; i++ {
			peerNode := "peer" + strconv.Itoa(i) + "." + orgHost
			peer := strings.Replace(entity, node, peerNode, -1)
			peers += peer
		}

	}
	orderers := `   
  orderer:`
	for i := 0; i < c.OrderCount; i++ {
		orderNode := "orderer" + strconv.Itoa(i) + "." + c.GetHostDomain(constant.OrdererSuffix)
		orderer := strings.Replace(entity, node, orderNode, -1)
		orderers += orderer
	}

	cas := `
  certificateAuthority:`
	for _, o := range c.PeersOrgs {
		caNode := "ca." + c.GetHostDomain(o)
		ca := strings.Replace(entity, node, caNode, -1)
		cas += ca
	}

	entityMatchers = entityMatchers + peers + orderers + cas
	c.buffer.WriteString(entityMatchers)
}

//node
func setNode(ntype, host, peerNode, crypto, notls string) string {
	peer := `
  {{node}}:
    url: {{node}}
    grpcOptions:
      ssl-target-name-override: {{node}}
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false
      allow-insecure: {{no-tls}}
    tlsCACerts:
      path: {{crypto-kubeconfig}}/{{type}}Organizations/{{domain}}/tlsca/tlsca.{{domain}}-cert.pem
`
	peer = strings.Replace(peer, domain, host, -1)
	peer = strings.Replace(peer, nodeType, ntype, -1)
	peer = strings.Replace(peer, node, peerNode, -1)
	peer = strings.Replace(peer, cryptoConfig, crypto, -1)

	if notls == "true" {
		notls = "false"
	} else {
		notls = "true"
	}
	peer = strings.Replace(peer, noTls, notls, -1)
	return peer
}

//构建
func (c *ConnectConfig) Build() []byte {
	c.setClient()
	c.setChannels()
	c.setOrganizations()
	c.setOrderers()
	c.setPeers()
	c.setCertificateAuthorities()
	c.setEntityMatchers()
	//写入
	ioutil.WriteFile(c.configYaml, c.buffer.Bytes(), os.ModePerm)

	return c.buffer.Bytes()
}

//获取字节流
func (c *ConnectConfig) GetBytes(maps map[string]interface{}) []byte {

	jsonStr, err := json.Marshal(maps)
	if err != nil {
		logger.Errorf("error: %v", err)
	}

	domains := new(model.ChainDomain)
	err = json.Unmarshal(jsonStr, domains)
	if err != nil {
		logger.Errorf("error: %v", err)
	}

	bts, err := ioutil.ReadFile(c.configYaml)
	if err != nil {
		logger.Errorf("error: %v", err)
	}

	for k, v := range domains.NodePorts {
		bts = bytes.Replace(bts, []byte(fmt.Sprintf(constant.Tag, k)), []byte(domains.NodeIps[0]+":"+v), -1)
	}
	return bts
}

func (c *ConnectConfig) setFabricChain(chain model.FabricChain) {
	c.ChainName = chain.ChainName
	c.Account = chain.Account
	c.Consensus = chain.Consensus
	c.PeersOrgs = chain.PeersOrgs
	c.OrderCount = chain.OrderCount
	c.PeerCount = chain.PeerCount
	c.ChannelName = chain.ChannelName
	c.TlsEnabled = chain.TlsEnabled
}

func NewConnectConfig(chain model.FabricChain, rootPath string) *ConnectConfig {
	config := &ConnectConfig{
		cryptoConfig: filepath.Join(rootPath, constant.CryptoConfigDir),
		buffer:       new(bytes.Buffer),
	}
	config.setFabricChain(chain)
	//链接文件
	config.configYaml = filepath.Join(rootPath, config.ChannelName+".yaml")
	return config
}
