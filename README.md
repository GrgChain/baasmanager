# baasmanager
![](https://img.shields.io/badge/build-passing-brightgreen.svg)
![](https://img.shields.io/badge/author-jonluo-yellow.svg)
![](https://img.shields.io/badge/kubernetes-v1.14.1-blue.svg)
![](https://img.shields.io/badge/go-v1.12.5-blue.svg)
![](https://img.shields.io/badge/docker-v18.06.3&ndash;ce-blue.svg)
![](https://img.shields.io/badge/hyperledger&nbsp;fabric-v1.4.1-blue.svg)

### 基于K8S平台的区块链即服务（Blockchain as a Service） 
### 整体功能
#### 动态创建fabric
- [x] solo
- [x] kafka
- [x] etcdraft
#### 区块链监控
- [x] 区块链首页统计分析 
- [ ] 区块链浏览器 
#### 区块链资源
- [x] 动态扩容
- [x] 释放 
### 主要目录结构
* baas-kubecluster  
  k8s集群，基于flannel网络，安装dashboard插件，还有其余插件等 (一个简单的k8s集群)
* baas-nfsshared  
  其会生成baas-artifacts，baas-fabric-data，baas-k8s-config目录  
  * baas-artifacts为存放生成的证书文件
  * baas-fabric-data为fabric网络映射出来的数据
  * baas-k8s-config为生成的k8s yaml定义文件  
* baas-template  
  fabric k8s的模板文件，用于生成baas-nfsshared/baas-k8s-config下的文件  
* baas-fabricengine  
  用于生成 baas-nfsshared的文件即目录结构和执行fabric操作
* baas-kubeengine  
  kubeconfig/config文件是k8s master的$HOME/.kube/config文件，用于k8s client链接k8s集群，将baas-nfsshared/baas-k8s-config下的文件在k8s集群创建启动  
* baas-gateway  
  统一api网关管理，调用入口
* baas-frontend  
  baas admin 前端
### 架构图
![](baas-others/images/baas.png)
### 数据流图
![](baas-others/images/flow.png)
### 页面
![](baas-others/images/das.png)
![](baas-others/images/user.png)
![](baas-others/images/role.png)
![](baas-others/images/chain.png)
![](baas-others/images/channel.png)
![](baas-others/images/chaincode.png)
![](baas-others/images/cc1.png)
![](baas-others/images/cc2.png)
### 部署样例
* [简单部署样例](sample.md)