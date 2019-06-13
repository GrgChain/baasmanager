### k8s engine
#### 用于管理 k8s
步骤
* 将k8s master的$HOME/.kube/config文件 替换 kubeconfig/config
* go get gitee.com/jonluo/baasmanager
* cd $GOPATH/src/gitee.com/jonluo/baasmanager/k8s-engine
* go run main.go --kubeconfig  $GOPATH/src/gitee.com/jonluo/baasmanager/k8s-engine/config/config