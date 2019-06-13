package engine

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type K8sEngine struct {
	jobs   *ItemQueue
	client *Clients
}

func NewK8sEngine(jobs *ItemQueue, client *Clients) *K8sEngine {
	return &K8sEngine{
		jobs:   jobs,
		client: client,
	}
}

func (k *K8sEngine) DoTasks() {

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
