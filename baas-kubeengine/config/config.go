package config

import (
	"github.com/spf13/viper"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
	"github.com/fsnotify/fsnotify"
	"os"
)

var Config *viper.Viper

var logger = log.GetLogger("kubuengine.config", log.INFO)

func init() {
	//监听改变动态跟新配置
	go watchConfig()
	//加载配置
	loadConfig()
}

//监听配置改变
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Config file changed:", e.Name)
		//改变重新加载
		loadConfig()
	})
}

//加载配置
func loadConfig() {
	viper.SetConfigName("keconfig")  // name of kubeconfig file
	viper.AddConfigPath(".")         // optionally look for kubeconfig in the working directory
	viper.AddConfigPath("/etc/baas") // path to look for the kubeconfig file in
	err := viper.ReadInConfig()      // Find and read the feconfig.yaml file
	if err != nil { // Handle errors reading the kubeconfig file
		logger.Errorf("Fatal error kubeconfig file: %s \n", err)
		os.Exit(-1)
	}
	//全局配置
	Config = viper.GetViper()
	logger.Infof("%v", Config.AllSettings())
}
