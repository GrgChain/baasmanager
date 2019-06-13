package engine

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"bytes"
	"io"
	"log"
)

var (
	kubeconfig   string
	apiServerURL string
)

type Clients struct {
	KubeClient *kubernetes.Clientset
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a config. Only required if out-of-cluster.")
	flag.StringVar(&apiServerURL, "master", "", "(Deprecated: switch to `--config`) The address of the Kubernetes API server. Overrides any value in config. Only required if out-of-cluster.")
}

// loadConfig loads a REST Config as per the rules specified in GetConfig
func loadConfig() (*rest.Config, error) {
	// If a flag is specified with the config location, use that
	if len(kubeconfig) > 0 {
		return clientcmd.BuildConfigFromFlags(apiServerURL, kubeconfig)
	}
	// If an env variable is specified with the config locaiton, use that
	if len(os.Getenv("KUBECONFIG")) > 0 {
		return clientcmd.BuildConfigFromFlags(apiServerURL, os.Getenv("KUBECONFIG"))
	}
	// If no explicit location, try the in-cluster config
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	// If no in-cluster config, try the default location in the user's home directory
	if usr, err := user.Current(); err == nil {
		if c, err := clientcmd.BuildConfigFromFlags(
			"", filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not locate a config")
}

func NewClients() *Clients {
	flag.Parse()
	// uses the current context in config
	config, err := loadConfig()
	if err != nil {
		log.Printf("Error building config: %s \n", err.Error())
	}

	if config.ServerName == "" {
		log.Printf("The cluster server name is %s \n", config.Host)
	} else {
		log.Printf("The cluster server name is %s \n", config.ServerName)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error building kubernetes clientset: %s \n", err.Error())
	}
	return &Clients{
		KubeClient: kubeClient,
	}

}

func (c *Clients) GetNodeList(ops metav1.ListOptions) *corev1.NodeList {

	ns, err := c.KubeClient.CoreV1().Nodes().List(ops)
	if err != nil {
		log.Printf(err.Error())
	}
	for _, n := range ns.Items {
		log.Println("Node：", n.Name, n.Status.Addresses)
	}
	return ns
}

func (c *Clients) GetNamespaceList(ops metav1.ListOptions) *corev1.NamespaceList {

	nss, err := c.KubeClient.CoreV1().Namespaces().List(ops)
	if err != nil {
		log.Printf(err.Error())
	}
	for _, ns := range nss.Items {
		log.Println("Namespace：", ns.Name, ns.Status.Phase)
	}
	return nss
}

func (c *Clients) GetPodList(ns string, ops metav1.ListOptions) *corev1.PodList {

	pods, err := c.KubeClient.CoreV1().Pods(ns).List(ops)
	if err != nil {
		log.Printf(err.Error())
	}
	for _, pod := range pods.Items {
		log.Println("Pod：", pod.Name, pod.Status.PodIP)
	}
	return pods
}

func (c *Clients) GetServiceList(ns string, ops metav1.ListOptions) *corev1.ServiceList {

	services, err := c.KubeClient.CoreV1().Services(ns).List(ops)
	if err != nil {
		log.Printf(err.Error())
	}
	for _, service := range services.Items {
		log.Println("Service：", service.Name, service.GetUID())
	}
	return services
}

func (c *Clients) CreateNameSpace(ns *corev1.Namespace) *corev1.Namespace {
	nameSpace, err := c.KubeClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created namesapce %q \n", nameSpace.GetObjectMeta().GetName())
	return nameSpace
}

func (c *Clients) CreatePod(pod *corev1.Pod) *corev1.Pod {

	newPod, err := c.KubeClient.CoreV1().Pods(pod.Namespace).Create(pod)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created pod %q \n", newPod.GetObjectMeta().GetName())
	return newPod
}

func (c *Clients) DeletePod(ns, name string, ops *metav1.DeleteOptions) {
	err := c.KubeClient.CoreV1().Pods(ns).Delete(name, ops)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (c *Clients) CreateDeployment(dep *appsv1.Deployment) *appsv1.Deployment {
	if dep.Namespace == "" {
		dep.Namespace = corev1.NamespaceDefault
	}
	deploymentsClient := c.KubeClient.AppsV1().Deployments(dep.Namespace)
	newDep, err := deploymentsClient.Create(dep)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created deployment %q \n", newDep.GetObjectMeta().GetName())
	return newDep
}

func (c *Clients) GetDeployment(ns string, depName string, ops metav1.GetOptions) *appsv1.Deployment {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(ns)
	redep, err := deploymentsClient.Get(depName, ops)
	if err != nil {
		log.Printf(err.Error())
	}
	return redep
}

func (c *Clients) GetDeploymentList(ns string, ops metav1.ListOptions) *appsv1.DeploymentList {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(ns)
	list, err := deploymentsClient.List(ops)
	if err != nil {
		log.Printf(err.Error())
	}
	for _, d := range list.Items {
		log.Println("Deployment ：", d.Name, d.Spec.Replicas)
	}
	return list
}

func (c *Clients) DeleteDeployment(ns string, depName string, ops *metav1.DeleteOptions) {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(ns)
	err := deploymentsClient.Delete(depName, ops)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (c *Clients) UpdateDeployment(ns string, depName string, dep *appsv1.Deployment) *appsv1.Deployment {
	deploymentsClient := c.KubeClient.AppsV1().Deployments(ns)

	newDep, err := deploymentsClient.Update(dep)
	if err != nil {
		log.Printf(err.Error())
	}

	return newDep
}

func (c *Clients) CreatePersistentVolume(pv *corev1.PersistentVolume) *corev1.PersistentVolume {

	newpv, err := c.KubeClient.CoreV1().PersistentVolumes().Create(pv)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created PersistentVolume %q \n", newpv.GetObjectMeta().GetName())
	return newpv
}

func (c *Clients) CreatePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) *corev1.PersistentVolumeClaim {
	newpvc, err := c.KubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(pvc)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created PersistentVolumeClaim %q \n", newpvc.GetObjectMeta().GetName())
	return newpvc
}

func (c *Clients) CreateService(service *corev1.Service) *corev1.Service {
	newservice, err := c.KubeClient.CoreV1().Services(service.Namespace).Create(service)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created Service %q \n", newservice.GetObjectMeta().GetName())
	return newservice
}

func (c *Clients) CreateStatefulSet(sfs *appsv1.StatefulSet) *appsv1.StatefulSet {
	if sfs.Namespace == "" {
		sfs.Namespace = corev1.NamespaceDefault
	}
	sfsClient := c.KubeClient.AppsV1().StatefulSets(sfs.Namespace)
	newSfs, err := sfsClient.Create(sfs)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Created deployment %q \n", newSfs.GetObjectMeta().GetName())
	return newSfs
}

func (c *Clients) PrintPodLogs(pod corev1.Pod) {
	podLogOpts := corev1.PodLogOptions{}

	req := c.KubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream()
	if err != nil {
		log.Println("error in opening stream")
	}
	if podLogs == nil {
		log.Println("error in opening stream")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		log.Println("error in copy information from podLogs to buf")
	}
	str := buf.String()

	log.Println("Pod logs :", str)
}
