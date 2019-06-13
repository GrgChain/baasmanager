### fabric engine
#### 用于生成fabric所需文件和执行fabric操作
步骤
* 与nfs server 部署在同一台机器
* 配置环境变量(相应修改ip)
  ```
  BaasRootPath=/home/jonluo/gopath/src/gitee.com/jonluo/baasmanager
  BaasNfsServer=192.168.1.45
  BaasK8sEngine=http://localhost:5991
  ```
* go get gitee.com/jonluo/baasmanager
* cd $GOPATH/src/gitee.com/jonluo/baasmanager/fabric-engine
* go run main.go 

