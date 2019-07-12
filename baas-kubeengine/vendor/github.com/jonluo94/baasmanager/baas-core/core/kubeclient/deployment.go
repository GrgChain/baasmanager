package kubeclient

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Clients) CreateDeployment(dep *appsv1.Deployment) *appsv1.Deployment {
	if dep.Namespace == "" {
		dep.Namespace = corev1.NamespaceDefault
	}
	deploymentsClient := c.KubeClient.AppsV1().Deployments(dep.Namespace)
	newDep, err := deploymentsClient.Create(dep)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Created deployment %q \n", newDep.GetObjectMeta().GetName())
	return newDep
}

func (c *Clients) GetDeployment(dep *appsv1.Deployment, ops metav1.GetOptions) *appsv1.Deployment {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(dep.Namespace)
	redep, err := deploymentsClient.Get(dep.Name, ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Get deployment %q \n", dep.GetObjectMeta().GetName())
	return redep
}

func (c *Clients) GetDeploymentList(dep *appsv1.Deployment, ops metav1.ListOptions) *appsv1.DeploymentList {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(dep.Namespace)
	list, err := deploymentsClient.List(ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	for _, d := range list.Items {
		logger.Infof("Deployment ：", d.Name, d.Spec.Replicas)
	}
	return list
}

func (c *Clients) DeleteDeployment(dep *appsv1.Deployment, ops *metav1.DeleteOptions) error {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(dep.Namespace)
	err := deploymentsClient.Delete(dep.Name, ops)
	if err != nil {
		logger.Errorf(err.Error())
	}
	logger.Infof("Delete deployment %q \n", dep.Name)
	return err
}

func (c *Clients) UpdateDeployment(dep *appsv1.Deployment) *appsv1.Deployment {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(dep.Namespace)

	newDep, err := deploymentsClient.Update(dep)
	if err != nil {
		logger.Errorf(err.Error())
	}

	logger.Infof("Updated deployment %q \n", newDep.GetObjectMeta().GetName())
	return newDep
}
