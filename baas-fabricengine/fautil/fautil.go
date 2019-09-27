package fautil

import (
	"strings"
	"github.com/jonluo94/baasmanager/baas-fabricengine/constant"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
	"path/filepath"
	"github.com/jonluo94/baasmanager/baas-core/common/util"
	"os"
	"github.com/jonluo94/baasmanager/baas-fabricengine/config"
)

func GetNamesapces(f model.FabricChain) string {
	namespaces := make([]string, 0)
	orderer := f.GetHostDomain(constant.OrdererSuffix)
	namespaces = append(namespaces, orderer)
	for _, o := range f.PeersOrgs {
		peer := f.GetHostDomain(o)
		namespaces = append(namespaces, peer)
	}
	return strings.Join(namespaces, ",")

}

func GetFirstOrderer(f model.FabricChain) string {
	return "orderer0." + f.GetHostDomain(constant.OrdererSuffix)
}

func GetChannelTx(f model.FabricChain, artifactPath string) string {
	return filepath.Join(artifactPath, constant.ChannelArtifactsDir, f.ChannelName+constant.Tx)
}

func GetAnchorsTxs(f model.FabricChain, artifactPath string) []string {
	anchorsTx := make([]string, len(f.PeersOrgs))
	for i, v := range f.PeersOrgs {
		tx := filepath.Join(artifactPath, constant.ChannelArtifactsDir, f.ChannelName+util.FirstUpper(v)+constant.AnchorsTx)
		anchorsTx[i] = tx
	}
	return anchorsTx
}

func GetChaincodeGithub(f model.FabricChannel) string {
	return filepath.Join(config.Config.GetString("BaasChaincodeGithub"), f.Account, f.ChainName,f.ChannelName, f.ChaincodeId, f.Version)
}

func GetChaincodeLocalGithub(f model.FabricChannel) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", config.Config.GetString("BaasChaincodeGithub"), f.Account, f.ChainName, f.ChannelName, f.ChaincodeId, f.Version)
}

func GetChaincodeGithubFile(f model.FabricChannel) string {
	return filepath.Join(GetChaincodeLocalGithub(f), constant.BaasChaincodeFile)
}
