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
		NodeIps:   nodeIps,
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
	for _, ns := range namesapces {
		podList := k.client.GetPodList(ns, metav1.ListOptions{})
		serviceList := k.client.GetServiceList(ns, metav1.ListOptions{})

		for _, pod := range podList.Items {
			name := ""
			podPort := int32(0)
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

			for _, ser := range serviceList.Items {
				if podPort != 0 {
					break
				}

				domain := ser.GetName() + "." + ser.GetNamespace()
				if domain != name {
					continue
				}

				for _, port := range ser.Spec.Ports {
					if port.Name == PortName {
						podPort = port.NodePort
						break
					}
				}
			}

			cp := model.ChainPod{
				Status:    string(pod.Status.Phase),
				HostIP:    string(pod.Status.HostIP),
				CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
				Name:      name,
				Port:      podPort,
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
