module github.com/jonluo94/baasmanager/baas-kubeengine

go 1.12

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-contrib/sessions v0.0.0-20190512062852-3cb4c4f2d615 // indirect
	github.com/gin-gonic/gin v1.7.0
	github.com/jonluo94/baasmanager/baas-core v0.0.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/spf13/viper v1.4.0
	k8s.io/api v0.0.0-20190703205437-39734b2a72fe
	k8s.io/apimachinery v0.0.0-20190703205208-4cfb76a8bf76
	k8s.io/client-go v12.0.0+incompatible // indirect
)

replace github.com/jonluo94/baasmanager/baas-core v0.0.0 => ../baas-core
