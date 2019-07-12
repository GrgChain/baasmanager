package kubeclient

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Clients) CreateStatefulSet(sfs *appsv1.StatefulSet) *appsv1.StatefulSet {
	if sfs.Namespace == "" {
		sfs.Namespace = corev1.NamespaceDefault
	}
	sfsClient := c.KubeClient.AppsV1().StatefulSets(sfs.Namespace)
	newSfs, err := sfsClient.Create(sfs)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Created deployment %q \n", newSfs.GetObjectMeta().GetName())
	return newSfs
}

func (c *Clients) DeleteStatefulSet(sfs *appsv1.StatefulSet, ops *metav1.DeleteOptions) error {
	if sfs.Namespace == "" {
		sfs.Namespace = corev1.NamespaceDefault
	}
	sfsClient := c.KubeClient.AppsV1().StatefulSets(sfs.Namespace)
	err := sfsClient.Delete(sfs.Name, ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Created deployment %q \n", sfs.GetObjectMeta().GetName())
	return err
}
