package kubeclient

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Clients) CreatePersistentVolume(pv *corev1.PersistentVolume) *corev1.PersistentVolume {

	newpv, err := c.KubeClient.CoreV1().PersistentVolumes().Create(pv)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Created PersistentVolume %q \n", newpv.GetObjectMeta().GetName())
	return newpv
}

func (c *Clients) DeletePersistentVolume(pv *corev1.PersistentVolume, ops *metav1.DeleteOptions) error {

	err := c.KubeClient.CoreV1().PersistentVolumes().Delete(pv.Name, ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Delete PersistentVolume %q \n", pv.GetObjectMeta().GetName())
	return err
}
