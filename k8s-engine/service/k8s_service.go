package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"gitee.com/jonluo/baasmanager/k8s-engine/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"gitee.com/jonluo/baasmanager/k8s-engine/engine"
)

const PortName = "endpoint"

type K8sService struct {
	client *engine.Clients
}

func NewK8sService(client *engine.Clients) *K8sService {
	return &K8sService{
		client: client,
	}
}

//部署
func (k *K8sService) deployData(ctx *gin.Context) {
	var pro engine.ParseProject
	if err := ctx.ShouldBindJSON(&pro); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	que := pro.Bytes2K8sEntities(pro.Yamls2Jsons())
	engine := engine.NewK8sEngine(que, k.client)
	engine.DoCreateTasks()
	models.ResultMsg(ctx, "success")
}


//删除
func (k *K8sService) deleteData(ctx *gin.Context) {
	var pro engine.ParseProject
	if err := ctx.ShouldBindJSON(&pro); err != nil {
		models.ResultFail(ctx, err)
		return
	}
	que := pro.Bytes2K8sEntities(pro.Yamls2Jsons())
	engine := engine.NewK8sEngine(que, k.client)
	engine.DoDeleteTasks()
	models.ResultMsg(ctx, "success")
}


func (k *K8sService) getChainDomain(ctx *gin.Context) {

	nss := ctx.Query("namesapces")
	namesapces := strings.Split(nss, ",")

	if len(namesapces) == 0 {
		models.ResultFail(ctx, "no namesapces")
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

	domains := models.NewChainDomain(nodeIps, portMap)
	models.ResultOk(ctx, domains)
}

func Server() {
	engine := engine.NewClients()
	k8sService := NewK8sService(engine)

	router := gin.Default()
	router.POST("/deployData", k8sService.deployData)
	router.POST("/deleteData", k8sService.deleteData)
	router.GET("/getChainDomain", k8sService.getChainDomain)
	router.Run(":5991")
}
