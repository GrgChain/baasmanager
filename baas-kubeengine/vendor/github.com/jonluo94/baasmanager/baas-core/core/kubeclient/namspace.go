package kubeclient

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Clients) GetNamespaceList(ops metav1.ListOptions) *corev1.NamespaceList {

	nss, err := c.KubeClient.CoreV1().Namespaces().List(ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	for _, ns := range nss.Items {
		logger.Infof("Namespace：", ns.Name, ns.Status.Phase)
	}
	return nss
}

func (c *Clients) CreateNameSpace(ns *corev1.Namespace) *corev1.Namespace {
	nameSpace, err := c.KubeClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Created namesapce %q \n", nameSpace.GetObjectMeta().GetName())
	return nameSpace
}

func (c *Clients) DeleteNameSpace(ns *corev1.Namespace, ops *metav1.DeleteOptions) error {
	err := c.KubeClient.CoreV1().Namespaces().Delete(ns.Name, ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Delete namesapce %q \n", ns.GetObjectMeta().GetName())
	return err
}
