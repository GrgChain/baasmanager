package generate

import (
	"io/ioutil"
	"log"
	"bytes"
	"os"
	"gitee.com/jonluo/baasmanager/fabric-engine/models"
	"gitee.com/jonluo/baasmanager/fabric-engine/constant"
	"strings"
	"gitee.com/jonluo/baasmanager/fabric-engine/config"
	"strconv"
	"fmt"
	"gitee.com/jonluo/baasmanager/fabric-engine/util"
	"path/filepath"
)

var (
	nfsArtifactPath   = fmt.Sprintf(constant.Tag, "nfs-artifact-path")
	nfsFabricDataPath = fmt.Sprintf(constant.Tag, "nfs-fabric-data-path")
	nfsServer         = fmt.Sprintf(constant.Tag, "nfs-server")
	namespace         = fmt.Sprintf(constant.Tag, "namespace")
	mspid             = fmt.Sprintf(constant.Tag, "mspid")
	index             = fmt.Sprintf(constant.Tag, "index")
	caPrivateKey      = fmt.Sprintf(constant.Tag, "ca-private-key")
	tlsEnabled        = fmt.Sprintf(constant.Tag, "tls-enabled")
)

type FabricK8s struct {
	ordererNamespace  string
	kafkaNamespace    string
	peerOrgsNamespace []string
	peerOrgsMsp       []string
	ordererCount      int
	peerCount         int
	consensus         string
	tlsEnabled        string
	nfsServer         string
	nfsActifactsPath  string
	nfsFabricDataPath string
	k8sYamlPath       string
	templatePath      string
}

func (f FabricK8s) buildNfsYaml() {
	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sNfsYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(nfsArtifactPath), []byte(f.nfsActifactsPath), -1)
	baseNs = bytes.Replace(baseNs, []byte(nfsFabricDataPath), []byte(f.nfsFabricDataPath), -1)
	baseNs = bytes.Replace(baseNs, []byte(nfsServer), []byte(f.nfsServer), -1)

	orderNs := bytes.Replace(baseNs, []byte(namespace), []byte(f.ordererNamespace), -1)
	buf.Write(orderNs)

	for _, v := range f.peerOrgsNamespace {
		peerNs := bytes.Replace(baseNs, []byte(namespace), []byte(v), -1)
		buf.Write(peerNs)
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sNfsYaml), buf.Bytes(), os.ModePerm)
}

func (f FabricK8s) buildNamespaceYaml() {

	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sNamespaceYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}

	orderNs := bytes.Replace(baseNs, []byte(namespace), []byte(f.ordererNamespace), -1)
	buf.Write(orderNs)

	for _, v := range f.peerOrgsNamespace {
		peerNs := bytes.Replace(baseNs, []byte(namespace), []byte(v), -1)
		buf.Write(peerNs)
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sNamespaceYaml), buf.Bytes(), os.ModePerm)

}

func (f FabricK8s) buildOrdererYaml() {

	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sOrdererYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(namespace), []byte(f.ordererNamespace), -1)
	baseNs = bytes.Replace(baseNs, []byte(tlsEnabled), []byte(f.tlsEnabled), -1)
	baseNs = bytes.Replace(baseNs, []byte(mspid), []byte(constant.OrdererMsp), -1)

	for i := 0; i < f.ordererCount; i++ {
		orderNs := bytes.Replace(baseNs, []byte(index), []byte(strconv.Itoa(i)), -1)
		buf.Write(orderNs)
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sOrdererYaml), buf.Bytes(), os.ModePerm)

}

func (f FabricK8s) buildPeerOrgYaml() {
	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sPeerYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(tlsEnabled), []byte(f.tlsEnabled), -1)

	for j, n := range f.peerOrgsNamespace {

		orgNs := bytes.Replace(baseNs, []byte(namespace), []byte(n), -1)
		orgNs = bytes.Replace(orgNs, []byte(mspid), []byte(f.peerOrgsMsp[j]), -1)

		for i := 0; i < f.peerCount; i++ {
			peerNs := bytes.Replace(orgNs, []byte(index), []byte(strconv.Itoa(i)), -1)
			buf.Write(peerNs)
		}
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sPeerYaml), buf.Bytes(), os.ModePerm)

}

func (f FabricK8s) buildCaYaml() {
	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sCaYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(tlsEnabled), []byte(f.tlsEnabled), -1)

	for _, n := range f.peerOrgsNamespace {
		caNs := bytes.Replace(baseNs, []byte(namespace), []byte(n), -1)

		caPath := filepath.Join(f.nfsActifactsPath, constant.CryptoConfigDir, "peerOrganizations", n, "ca")
		var caPriKey string
		files, _ := ioutil.ReadDir(caPath)
		for _, f := range files {
			name := f.Name()
			if strings.HasSuffix(name, constant.PriKeySuf) {
				caPriKey = name
			}
		}

		caNs = bytes.Replace(caNs, []byte(caPrivateKey), []byte(caPriKey), -1)

		buf.Write(caNs)
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sCaYaml), buf.Bytes(), os.ModePerm)

}
func (f FabricK8s) buildCliYaml() {
	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sCliYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(tlsEnabled), []byte(f.tlsEnabled), -1)

	for j, n := range f.peerOrgsNamespace {
		caNs := bytes.Replace(baseNs, []byte(namespace), []byte(n), -1)
		caNs = bytes.Replace(caNs, []byte(mspid), []byte(f.peerOrgsMsp[j]), -1)
		buf.Write(caNs)
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sCliYaml), buf.Bytes(), os.ModePerm)

}
func (f FabricK8s) buildZookeeperYaml() {
	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sZookeeperYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(namespace), []byte(f.kafkaNamespace), -1)
	baseNs = bytes.Replace(baseNs, []byte(nfsFabricDataPath), []byte(f.nfsFabricDataPath), -1)
	baseNs = bytes.Replace(baseNs, []byte(nfsServer), []byte(f.nfsServer), -1)
	buf.Write(baseNs)
	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sZookeeperYaml), buf.Bytes(), os.ModePerm)

}
func (f FabricK8s) buildKafkaYaml() {
	var buf bytes.Buffer
	baseNs, err := ioutil.ReadFile(filepath.Join(f.templatePath, constant.K8sKafkaYaml))
	if err != nil {
		log.Printf("error: %v", err)
	}
	baseNs = bytes.Replace(baseNs, []byte(namespace), []byte(f.kafkaNamespace), -1)
	for i := 0; i < 4; i++ {
		brokerNs := bytes.Replace(baseNs, []byte(index), []byte(strconv.Itoa(i)), -1)
		buf.Write(brokerNs)
	}

	ioutil.WriteFile(filepath.Join(f.k8sYamlPath, constant.K8sKafkaYaml), buf.Bytes(), os.ModePerm)

}

func (f FabricK8s) Build() {
	f.buildNamespaceYaml()
	f.buildNfsYaml()
	f.buildOrdererYaml()
	f.buildPeerOrgYaml()
	f.buildCaYaml()
	f.buildCliYaml()

	switch f.consensus {
	case constant.OrdererSolo:
		log.Println("solo")
	case constant.OrdererKafka:
		f.buildZookeeperYaml()
		f.buildKafkaYaml()
	}
}

func NewFabricK8s(chain models.FabricChain, conf config.UserBaasConfig) FabricK8s {

	peerDamoins := make([]string, len(chain.PeersOrgs))
	peerMsps := make([]string, len(chain.PeersOrgs))

	for i, v := range chain.PeersOrgs {
		peer := chain.GetHostDomain(v)
		peerDamoins[i] = peer

		msp := util.FirstUpper(v) + constant.MspSuf
		peerMsps[i] = msp
	}

	return FabricK8s{
		ordererNamespace:  chain.GetHostDomain(constant.OrdererSuffix),
		kafkaNamespace:    chain.GetHostDomain(constant.KafkaSuffix),
		peerOrgsNamespace: peerDamoins,
		peerOrgsMsp:       peerMsps,
		ordererCount:      chain.OrderCount,
		peerCount:         chain.PeerCount,
		consensus:         chain.Consensus,
		tlsEnabled:        chain.TlsEnabled,
		nfsServer:         constant.BaasNfsServer,
		nfsActifactsPath:  conf.ArtifactPath,
		nfsFabricDataPath: conf.DataPath,
		k8sYamlPath:       conf.K8sConfigPath,
		templatePath:      conf.TemplatePath,
	}
}
