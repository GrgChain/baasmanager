package engine

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/jonluo94/baasmanager/baas-core/core/kubeclient"
	"github.com/jonluo94/baasmanager/baas-core/common/queue"
	"github.com/jonluo94/baasmanager/baas-core/common/json"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
)

var logger = log.GetLogger("kubeengine.engine", log.INFO)

const (
	ApiVersion        = "v1"
	DeploymentVersion = "apps/v1beta1"

	KindDeployment            = "Deployment"
	KindNamespace             = "Namespace"
	KindService               = "Service"
	KindPersistentVolume      = "PersistentVolume"
	KindPersistentVolumeClaim = "PersistentVolumeClaim"
	KindStatefulSet           = "StatefulSet"
)

type KubeEngine struct {
	jobs   *queue.Queue
	client *kubeclient.Clients
}

func NewKubeEngine(jobs *queue.Queue, client *kubeclient.Clients) *KubeEngine {
	return &KubeEngine{
		jobs:   jobs,
		client: client,
	}
}

func (k *KubeEngine) DoCreateTasks() {

	for !k.jobs.IsEmpty() {
		item := k.jobs.Dequeue()
		switch item.(type) {
		case *corev1.Namespace:
			k.client.CreateNameSpace(item.(*corev1.Namespace))
		case *appsv1.Deployment:
			k.client.CreateDeployment(item.(*appsv1.Deployment))
		case *corev1.Service:
			k.client.CreateService(item.(*corev1.Service))
		case *corev1.PersistentVolume:
			k.client.CreatePersistentVolume(item.(*corev1.PersistentVolume))
		case *corev1.PersistentVolumeClaim:
			k.client.CreatePersistentVolumeClaim(item.(*corev1.PersistentVolumeClaim))
		case *appsv1.StatefulSet:
			k.client.CreateStatefulSet(item.(*appsv1.StatefulSet))
		}

	}
}

func (k *KubeEngine) DoDeleteTasks() {

	for !k.jobs.IsEmpty() {
		item := k.jobs.Dequeue()
		deletePolicy := metav1.DeletePropagationForeground
		delops := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
			switch item.(type) {
		case *corev1.Namespace:
			k.client.DeleteNameSpace(item.(*corev1.Namespace),delops)
		case *appsv1.Deployment:
			k.client.DeleteDeployment(item.(*appsv1.Deployment),delops)
		case *corev1.Service:
			k.client.DeleteService(item.(*corev1.Service),delops)
		case *corev1.PersistentVolume:
			k.client.DeletePersistentVolume(item.(*corev1.PersistentVolume),delops)
		case *corev1.PersistentVolumeClaim:
			k.client.DeletePersistentVolumeClaim(item.(*corev1.PersistentVolumeClaim),delops)
		case *appsv1.StatefulSet:
			k.client.DeleteStatefulSet(item.(*appsv1.StatefulSet),delops)
		}

	}
}


func Bytes2K8sEntities(jsonArray [][]byte) *queue.Queue {

	queue := queue.NewQueue()

	for _, jsonObj := range jsonArray {
		var entity map[string]interface{}
		err := json.Unmarshal(jsonObj, &entity)
		if err != nil {
			logger.Error(err.Error())
		}
		if entity["apiVersion"] == ApiVersion || entity["apiVersion"] == DeploymentVersion {
			switch entity["kind"] {
			case KindNamespace:
				namespace := &corev1.Namespace{}
				err = json.Unmarshal(jsonObj, &namespace)
				if err != nil {
					logger.Error(err.Error())
				}
				queue.Enqueue(namespace)
			case KindDeployment:
				deployment := &appsv1.Deployment{}
				err = json.Unmarshal(jsonObj, &deployment)
				if err != nil {
					logger.Error(err.Error())
				}
				queue.Enqueue(deployment)
			case KindService:
				service := &corev1.Service{}
				err = json.Unmarshal(jsonObj, &service)
				if err != nil {
					logger.Error(err.Error())
				}
				queue.Enqueue(service)
			case KindPersistentVolume:
				pv := &corev1.PersistentVolume{}
				err = json.Unmarshal(jsonObj, &pv)
				if err != nil {
					logger.Error(err.Error())
				}
				queue.Enqueue(pv)
			case KindPersistentVolumeClaim:
				pvc := &corev1.PersistentVolumeClaim{}
				err = json.Unmarshal(jsonObj, &pvc)
				if err != nil {
					logger.Error(err.Error())
				}
				queue.Enqueue(pvc)
			case KindStatefulSet:
				sfs := &appsv1.StatefulSet{}
				err = json.Unmarshal(jsonObj, &sfs)
				if err != nil {
					logger.Error(err.Error())
				}
				queue.Enqueue(sfs)
			}
		}
	}
	return queue
}
