module github.com/jonluo94/baasmanager/baas-fabricengine

go 1.12

require (
	github.com/Shopify/sarama v1.23.0 // indirect
	github.com/cloudflare/cfssl v0.0.0-20180323000720-5d63dbd981b5 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/fsouza/go-dockerclient v1.4.1 // indirect
	github.com/gin-contrib/sessions v0.0.0-20190512062852-3cb4c4f2d615 // indirect
	github.com/gin-gonic/gin v1.7.0
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hyperledger/fabric v1.4.1 // indirect
	github.com/hyperledger/fabric-amcl v0.0.0-20181230093703-5ccba6eab8d6 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0-alpha5.0.20190411180201-5a9a0e749e4f // indirect
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190411180201-5a9a0e749e4f // indirect
	github.com/jonluo94/baasmanager/baas-core v0.0.0
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mitchellh/mapstructure v0.0.0-20180511142126-bb74f1db0675 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/prometheus/common v0.0.0-20180801064454-c7de2306084e // indirect
	github.com/prometheus/procfs v0.0.0-20180920065004-418d78d0b9a7 // indirect
	github.com/spf13/afero v1.1.1 // indirect
	github.com/spf13/viper v1.0.2
	github.com/sykesm/zap-logfmt v0.0.2 // indirect
	go.uber.org/zap v1.10.0 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.0.0-20190703205208-4cfb76a8bf76 // indirect
)

replace github.com/jonluo94/baasmanager/baas-core v0.0.0 => ../baas-core
