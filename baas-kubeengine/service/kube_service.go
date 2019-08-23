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
	"strconv"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
)

const PortName = "endpoint"

var logger = log.GetLogger("kubeengine.service", log.INFO)

type KubeService struct {
	client *kubeclient.Clients
}

func NewKubeService(client *kubeclient.Clients) *KubeService {
	return &KubeService{
		client: client,
	}
}

//获取节点ip
func (k *KubeService) getNodeIPs() []string{

	nodeList := k.client.GetNodeList(metav1.ListOptions{})
	nodeIPs := make([]string, len(nodeList.Items))
	for i, node := range nodeList.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				nodeIPs[i] = addr.Address
				break
			}
		}
	}
	return nodeIPs
}
//获取服务map
func (k *KubeService) getServiceMap(namesapces []string) map[string]string{
	portMap := make(map[string]string, 0)
	for _, ns := range namesapces {
		//获取服务
		serviceList := k.client.GetServiceList(ns, metav1.ListOptions{})
		for _, ser := range serviceList.Items {
			//获取服务域名
			domain := ser.GetName() + "." + ser.GetNamespace()
			for _, port := range ser.Spec.Ports {
				if port.Name == PortName {
					portMap[domain] = fmt.Sprintf("%d", port.NodePort)
					break
				}
			}
		}
	}
	return portMap
}

//部署
func (k *KubeService) DeployData(ctx *gin.Context) {
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
func (k *KubeService) DeleteData(ctx *gin.Context) {
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

func (k *KubeService) GetChainDomain(ctx *gin.Context) {

	nss := ctx.Query("namesapces")
	namesapces := strings.Split(nss, ",")

	if len(namesapces) == 0 {
		gintool.ResultFail(ctx, "no namesapces")
	}
    //获取节点ip
    nondeIPs := k.getNodeIPs()
    //获取服务map
	portMap := k.getServiceMap(namesapces)

	domains := model.ChainDomain{
		NodeIps:   nondeIPs,
		NodePorts: portMap,
	}

	gintool.ResultOk(ctx, domains)
}




func (k *KubeService) GetChainPods(ctx *gin.Context) {

	nss := ctx.Query("namesapces")
	namesapces := strings.Split(nss, ",")

	if len(namesapces) == 0 {
		gintool.ResultFail(ctx, "no namesapces")
	}
	chainPods := make([]model.ChainPod, 0)

	portMap := k.getServiceMap(namesapces)

	for _, ns := range namesapces {
		podList := k.client.GetPodList(ns, metav1.ListOptions{})

		for _, pod := range podList.Items {
			name := ""
			podType := pod.Labels["role"]

			switch podType {
			case "ca":
				name = pod.Labels["name"] + "." + pod.Namespace
			case "orderer":
				name = pod.Labels["orderer-id"] + "." + pod.Namespace
			case "peer":
				name = pod.Labels["peer-id"] + "." + pod.Namespace
			case "kafka", "zookeeper":
				continue
			}

			podPort ,err := strconv.Atoi(portMap[name])
			if err != nil {
               logger.Error(err)
			}

			cp := model.ChainPod{
				Status:    string(pod.Status.Phase),
				HostIP:    string(pod.Status.HostIP),
				CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Name:      name,
				Port:      int32(podPort),
				Type:      podType,
			}
			chainPods = append(chainPods, cp)
		}

	}

	gintool.ResultOk(ctx, chainPods)
}

func Server() {
	kc := kubeclient.NewClients(config.Config.GetString("BaasKubeMasterConfig"))
	kubeService := NewKubeService(kc)

	router := gin.New()
	router.Use(gintool.Logger())
	router.Use(gin.Recovery())
	router.POST("/deployData", kubeService.DeployData)
	router.POST("/deleteData", kubeService.DeleteData)
	router.GET("/getChainDomain", kubeService.GetChainDomain)
	router.GET("/getChainPods", kubeService.GetChainPods)
	router.Run(":" + config.Config.GetString("BaasKubeEnginePort"))
}
