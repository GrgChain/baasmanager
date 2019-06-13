package service

import (
	"gitee.com/jonluo/baasmanager/fabric-engine/constant"
	"gitee.com/jonluo/baasmanager/fabric-engine/models"
	"gitee.com/jonluo/baasmanager/fabric-engine/util"
)

type K8sService struct {
}

func (k K8sService) deployData(datas *models.K8sData) []byte {
	return util.PostJson(constant.BaasK8sEngine+"/deployData", datas)
}

func (k K8sService) getChainDomain(nss string) []byte {
	return util.Get(constant.BaasK8sEngine + "/getChainDomain?namesapces=" + nss)
}
