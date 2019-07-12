package kubeclient

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Clients) CreatePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) *corev1.PersistentVolumeClaim {
	newpvc, err := c.KubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(pvc)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Created PersistentVolumeClaim %q \n", newpvc.GetObjectMeta().GetName())
	return newpvc
}

func (c *Clients) DeletePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, ops *metav1.DeleteOptions) error {
	err := c.KubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(pvc.Name, ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Delete PersistentVolumeClaim %q \n", pvc.GetObjectMeta().GetName())
	return err
}
