package kubeclient

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
)

var logger = log.GetLogger("kubeclient", log.INFO)

type Clients struct {
	KubeClient *kubernetes.Clientset
}

// loadConfig loads a REST Config as per the rules specified in GetConfig
func loadConfig(kubeconfig string) (*rest.Config, error) {
	//集群内部地址
	var apiServerURL = ""
	// If a flag is specified with the kubeconfig location, use that
	if len(kubeconfig) > 0 {
		return clientcmd.BuildConfigFromFlags(apiServerURL, kubeconfig)
	}
	// If an env variable is specified with the kubeconfig locaiton, use that
	if len(os.Getenv("KUBECONFIG")) > 0 {
		return clientcmd.BuildConfigFromFlags(apiServerURL, os.Getenv("KUBECONFIG"))
	}
	// If no explicit location, try the in-cluster kubeconfig
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	// If no in-cluster kubeconfig, try the default location in the user's home directory
	if usr, err := user.Current(); err == nil {
		if c, err := clientcmd.BuildConfigFromFlags(
			"", filepath.Join(usr.HomeDir, ".kube", "kubeconfig")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not locate a kubeconfig")
}

func NewClients(kubeconfig string) *Clients {
	// uses the current context in kubeconfig
	config, err := loadConfig(kubeconfig)
	if err != nil {
		logger.Errorf("Error building kubeconfig: %s \n", err.Error())
	}

	if config.ServerName == "" {
		logger.Infof("The cluster server name is %s \n", config.Host)
	} else {
		logger.Infof("The cluster server name is %s \n", config.ServerName)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Errorf("Error building kubernetes clientset: %s \n", err.Error())
	}
	return &Clients{
		KubeClient: kubeClient,
	}

}
