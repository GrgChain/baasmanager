## k8s 简单集群搭建
1 master 1 node 集群方便演示
* centos7
* k8s 1.14.1
* docker 18.06.3

两台centos7虚拟机

| 主机名 | ip |
|:----:|:----:|
| k8s-master | 192.168.0.224 | 
| k8s-node1 | 192.168.0.51 |   

### k8s-master 
通过ssh操作 (ssh root@192.168.0.224)

1.修改主机名
  ```
  hostnamectl --static set-hostname k8s-master 
  echo -e "192.168.0.224 k8s-master\n192.168.0.51 k8s-node1" >> /etc/hosts  
  ```
2.关闭防火墙
  ```
  systemctl stop firewalld & systemctl disable firewalld
  ```
3.关闭SeLinux
  ```
  setenforce 0
  ```
4.关闭Swap
  ```
  sed -i '/ swap / s/^/#/' /etc/fstab
  ```
5.安装docker  
  安装需要的软件包， yum-util 提供yum-config-manager功能，另外两个是devicemapper驱动依赖的
  ``` 
  yum install -y yum-utils device-mapper-persistent-data lvm2
  ```
  设置yum源
  ```
  yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo 
  ```
  可以查看所有仓库中所有docker版本，并选择特定版本安装
  ```
  yum list docker-ce --showduplicates | sort -r
  ```
  安装
  ```
  yum install docker-ce-18.06.3.ce-3.el7   
  ```
  启动
  ```
  systemctl enable docker && systemctl start docker  
  ```
  配置镜像加速器,通过修改daemon配置文件/etc/docker/daemon.json来使用加速器
  ``` 
  sudo mkdir -p /etc/docker
  sudo tee /etc/docker/daemon.json <<-'EOF'
  {
    "registry-mirrors": ["https://ijdk512y.mirror.aliyuncs.com"]
  }
  EOF
  sudo systemctl daemon-reload
  sudo systemctl restart docker
  ```
  验证
  ```
  docker version
  ```

6.安装kubelet、kubeadm、kubectl  
  配置K8S的yum源,执行以下命令安装kubelet、kubeadm、kubectl
  ``` 
  cat <<EOF > /etc/yum.repos.d/kubernetes.repo
  [kubernetes]
  name=Kubernetes
  baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
  enabled=1
  gpgcheck=0
  repo_gpgcheck=0
  gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
         http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
  EOF
  ```
  安装
  ```
  yum install -y kubelet-1.14.1-0.x86_64 kubeadm-1.14.1-0.x86_64 kubectl-1.14.1-0.x86_64
  ```
  CentOS 7上的一些用户报告了由于iptables被绕过而导致流量路由不正确的问题。应该确保 net.bridge.bridge-nf-call-iptables在sysctl配置中设置为1
  ```
  cat <<EOF >  /etc/sysctl.d/k8s.conf
  net.bridge.bridge-nf-call-ip6tables = 1
  net.bridge.bridge-nf-call-iptables = 1
  EOF
  sysctl --system
  ```
  启动k8s相关服务
  ``` 
  systemctl enable kubelet && systemctl start kubelet
  ```    

7.初始化k8s集群镜像准备，默认去访问谷歌的服务器，以下载集群所依赖的Docker镜像,因此也会超时失败。docker.io仓库对google的容器做了镜像，可以通过下列命令下拉取相关镜像 
  ```
  docker pull mirrorgooglecontainers/kube-apiserver-amd64:v1.14.1
  docker pull mirrorgooglecontainers/kube-controller-manager-amd64:v1.14.1
  docker pull mirrorgooglecontainers/kube-scheduler-amd64:v1.14.1
  docker pull mirrorgooglecontainers/kube-proxy-amd64:v1.14.1
  docker pull mirrorgooglecontainers/pause:3.1
  docker pull mirrorgooglecontainers/etcd-amd64:3.3.10
  docker pull coredns/coredns:1.3.1
  docker pull thejosan20/flannel:v0.10.0-amd64
  docker pull mirrorgooglecontainers/kubernetes-dashboard-amd64:v1.10.0
  
  ```
  版本信息需要根据实际情况进行相应的修改。通过docker tag命令来修改镜像的标签
  ``` 
  docker tag docker.io/mirrorgooglecontainers/kube-proxy-amd64:v1.14.1 k8s.gcr.io/kube-proxy:v1.14.1
  docker tag docker.io/mirrorgooglecontainers/kube-scheduler-amd64:v1.14.1 k8s.gcr.io/kube-scheduler:v1.14.1
  docker tag docker.io/mirrorgooglecontainers/kube-apiserver-amd64:v1.14.1 k8s.gcr.io/kube-apiserver:v1.14.1
  docker tag docker.io/mirrorgooglecontainers/kube-controller-manager-amd64:v1.14.1 k8s.gcr.io/kube-controller-manager:v1.14.1
  docker tag docker.io/mirrorgooglecontainers/etcd-amd64:3.3.10  k8s.gcr.io/etcd:3.3.10
  docker tag docker.io/mirrorgooglecontainers/pause:3.1  k8s.gcr.io/pause:3.1
  docker tag docker.io/coredns/coredns:1.3.1  k8s.gcr.io/coredns:1.3.1
  docker tag docker.io/thejosan20/flannel:v0.10.0-amd64 quay.io/coreos/flannel:v0.10.0-amd64
  docker tag docker.io/mirrorgooglecontainers/kubernetes-dashboard-amd64:v1.10.0  k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.0
  ```
  删除原镜像
  ```
  docker rmi mirrorgooglecontainers/kube-apiserver-amd64:v1.14.1
  docker rmi mirrorgooglecontainers/kube-controller-manager-amd64:v1.14.1
  docker rmi mirrorgooglecontainers/kube-scheduler-amd64:v1.14.1
  docker rmi mirrorgooglecontainers/kube-proxy-amd64:v1.14.1
  docker rmi mirrorgooglecontainers/pause:3.1
  docker rmi mirrorgooglecontainers/etcd-amd64:3.3.10
  docker rmi coredns/coredns:1.3.1
  docker rmi thejosan20/flannel:v0.10.0-amd64
  docker rmi mirrorgooglecontainers/kubernetes-dashboard-amd64:v1.10.0  
  ```
  fabric 相关镜像
  ```
  docker pull hyperledger/fabric-ca
  docker pull hyperledger/fabric-tools
  docker pull hyperledger/fabric-orderer
  docker pull hyperledger/fabric-peer
  docker pull hyperledger/fabric-javaenv
  docker pull hyperledger/fabric-ccenv
  docker pull hyperledger/fabric-zookeeper
  docker pull hyperledger/fabric-kafka
  docker pull hyperledger/fabric-couchdb
  ```

8.安装 nfs 工具后重启
  ``` 
  yum -y install nfs-utils 
  reboot
  ```
  
### k8s-node1 
通过ssh操作 (ssh root@192.168.0.51)

1.修改主机名
  ```
  hostnamectl --static set-hostname k8s-node1
  echo -e "192.168.0.224 k8s-master\n192.168.0.51 k8s-node1"  >> /etc/hosts  
  ```
其余步骤与k8s-master一样


## 环境准备好后，开始搭建集群
* 进入 k8s-master (ssh root@192.168.0.224)
  * kubeadm 初始化 master (pod-network-cidr flannel网络用到)
    ```
    kubeadm init --kubernetes-version=v1.14.1 --pod-network-cidr=10.244.0.0/16
    ```  
    执行日志中脚本
    ```
    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config
    ```  
* 进入 k8s-node1 (ssh root@192.168.0.51)
  ```
  kubeadm join 192.168.0.224:6443 --token l1tmiz.mugfz7lenp7ih28g \
      --discovery-token-ca-cert-hash sha256:d4cf165051ac0d45386bb6168429af8aa2cda504d884d1ff499a168c2a711209
  ```

* 进入 k8s-master (ssh root@192.168.0.224)
  * 在k8s-master验证，并创建网络
    ```
    kubectl get nodes
    kubectl taint nodes --all node-role.kubernetes.io/master-
    ```
    编辑flannel/kube-flannel.yml,创建flannel网络
    ```
    kubectl apply -f kube-flannel.yml
    ```
    查看pods
    ```
    kubectl get pods --all-namespaces
    ```
    编辑dashboard/kubernetes-dashboard.yaml,创建K8S Dashboard
    ```
    kubectl create -f kubernetes-dashboard.yaml
    ```
    编辑dashboard/admin-token.yaml,创建Dashboard 管理员用户
    ```
    kubectl create -f admin-token.yaml
    ```
    获取登陆token
    ```
    kubectl describe secret/$(kubectl get secret -nkube-system |grep admin|awk '{print $1}') -nkube-system
    ```
    浏览器打开：https://192.168.0.224:30000/#!/login 令牌为token登录
    
  * 获取 kube-dns 的ip地址
    ```
    kubectl get services --all-namespaces | grep kube-dns 
    ```
    log
    ``` 
    kube-system   kube-dns               ClusterIP   10.96.0.10      <none>        53/UDP,53/TCP,9153/TCP   26m
    ```
    得到kube-dns的ip：10.96.0.10

* 在k8s集群搭建完后操作  
  为了解决解析域名的问题，需要在k8s集群每个worker节点的 ExecStart 中加入相关参数：
  kube-dns 的 ip 为10.96.0.10，宿主机网络 DNS 的地址为 192.168.0.1，
  为使得 chaincode 的容器可以解析到 peer 节点，在每个k8s worker节点，修改步骤如下
  ``` 
  vi /lib/systemd/system/docker.service
  ```
  在 ExecStart 参数后追加：
  ```
  --dns=10.96.0.10 --dns=192.168.0.1 --dns-search default.svc.cluster.local --dns-search svc.cluster.local --dns-opt ndots:2 --dns-opt timeout:2 --dns-opt attempts:2
  ```
  重启docker
  ```
  systemctl daemon-reload && systemctl restart docker 
  ```

