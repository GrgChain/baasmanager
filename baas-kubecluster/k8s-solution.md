### K8S Dashboard  
```
kubectl create -f dashboard/kubernetes-dashboard.yaml 
```
浏览器 https://192.168.1.210:30000/#!/login 即可以验证dashboard是否安装成功.查看秘钥, 列出所有kube-system命名空间下秘钥
```
kubectl -n kube-system get secret
```
创建用户
```
kubectl create -f dashboard/admin-token.yaml
```
获取登陆token
```
kubectl describe secret/$(kubectl get secret -nkube-system |grep admin|awk '{print $1}') -nkube-system
```

### K8S监控方案 Prometheus Grafana
Kubernetes Setup for Prometheus and Grafana
https://github.com/giantswarm/kubernetes-prometheus
安装
```
kubectl create -f prometheus/kubernetes-prometheus.yaml
```

要再次关闭所有组件，只需删除该命名空间
```
kubectl delete namespace monitoring
```
### K8S日志方案 EFK

* 给 Node 设置标签
```
kubectl get nodes
```

```
kubectl label nodes k8s-master beta.kubernetes.io/fluentd-ds-ready=true
kubectl label nodes k8s-node1 beta.kubernetes.io/fluentd-ds-ready=true
```
``` 
kubectl create -f ./efk/
```

注意: kibana Pod 第一次启动时会用较长时间(10-20分钟)来优化和 Cache 状态页面

* 通过 kubectl proxy 访问：  
在master节点创建代理
```
kubectl proxy --port=8888 --address=192.168.1.210 --accept-hosts=^*$
```
浏览器访问 URL：
``` 
http://192.168.1.210:8888/api/v1/namespaces/kube-system/services/kibana-logging/proxy
```

### K8S 外网访问代理 nginx-ingress
pull镜像：  
```
docker pull siriuszg/nginx-ingress-controller:0.20.0  
docker tag docker.io/siriuszg/nginx-ingress-controller:0.20.0 quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.20.0
```
demo.yaml是 ingress应用样例
``` 
kubectl create -f ingress.yaml
kubectl create -f demo.yaml
```
demo中 registry.my.nginx 为映射 Host  
访问：  
curl -H "Host: registry.my.nginx" http://192.168.0.224:30080
