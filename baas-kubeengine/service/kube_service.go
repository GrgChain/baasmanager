package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
	"github.com/jonluo94/baasmanager/baas-kubeengine/engine"
	"github.com/jonluo94/baasmanager/baas-core/core/kubeclient"
	"github.com/jonluo94/baasmanager/baas-kubeengine/config"
	"github.com/jonluo94/baasmanager/baas-core/common/util"
)

const PortName = "endpoint"

type KubeService struct {
	client *kubeclient.Clients
}

func NewKubeService(client *kubeclient.Clients) *KubeService {
	return &KubeService{
		client: client,
	}
}

//部署
func (k *KubeService) deployData(ctx *gin.Context) {
	var pro model.K8sData
	if err := ctx.ShouldBindJSON(&pro); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	que := engine.Bytes2K8sEntities(util.Yamls2Jsons(pro.Data))
	engine := engine.NewKubeEngine(que, k.client)
	engine.DoCreateTasks()
	gintool.ResultMsg(ctx, "success")
}


//删除
func (k *KubeService) deleteData(ctx *gin.Context) {
	var pro model.K8sData
	if err := ctx.ShouldBindJSON(&pro); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	que := engine.Bytes2K8sEntities(util.Yamls2Jsons(pro.Data))
	engine := engine.NewKubeEngine(que, k.client)
	engine.DoDeleteTasks()
	gintool.ResultMsg(ctx, "success")
}


func (k *KubeService) getChainDomain(ctx *gin.Context) {

	nss := ctx.Query("namesapces")
	namesapces := strings.Split(nss, ",")

	if len(namesapces) == 0 {
		gintool.ResultFail(ctx, "no namesapces")
	}

	nodeIps := make([]string, 0)
	nodeList := k.client.GetNodeList(metav1.ListOptions{})
	for _, node := range nodeList.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				nodeIps = append(nodeIps, addr.Address)
				break
			}
		}
	}

	portMap := make(map[string]string, 0)
	for _, ns := range namesapces {
		serviceList := k.client.GetServiceList(ns, metav1.ListOptions{})
		for _, ser := range serviceList.Items {
			domain := ser.GetName() + "." + ser.GetNamespace()
			for _, port := range ser.Spec.Ports {
				if port.Name == PortName {
					portMap[domain] = fmt.Sprintf("%d", port.NodePort)
					break
				}
			}
		}
	}

	domains := model.ChainDomain{
		NodeIps:nodeIps,
		NodePorts:portMap,
	}

	gintool.ResultOk(ctx, domains)
}

func Server() {
	kc := kubeclient.NewClients(config.Config.GetString("BaasKubeMasterConfig"))
	kubeService := NewKubeService(kc)

	router := gin.New()
	router.Use(gintool.Logger())
	router.Use(gin.Recovery())
	router.POST("/deployData", kubeService.deployData)
	router.POST("/deleteData", kubeService.deleteData)
	router.GET("/getChainDomain", kubeService.getChainDomain)
	router.Run(":"+config.Config.GetString("BaasKubeEnginePort"))
}
