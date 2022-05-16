module github.com/jonluo94/baasmanager/baas-gateway

go 1.12

require (
	github.com/alexandrevicenzi/unchained v0.0.0-20190214114102-ecd422680cf1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-contrib/sessions v0.0.0-20190512062852-3cb4c4f2d615 // indirect
	github.com/gin-gonic/gin v1.7.0
	github.com/go-xorm/xorm v0.7.1
	github.com/jonluo94/baasmanager/baas-core v0.0.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/spf13/viper v1.4.0
)

replace github.com/jonluo94/baasmanager/baas-core v0.0.0 => ../baas-core
