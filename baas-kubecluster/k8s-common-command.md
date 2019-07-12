## K8S常用命令

#### 版本
* kubectl version

#### 集群构成
* kubectl get nodes

#### 运行一个镜像
* kubectl run nginx-dev --image=nginx:1.15 --replicas=1 --port=8080

#### 确认Deployment
* kubectl get deployment

#### 确认pod
* kubectl get pods

#### 删除pod
* kubectl delete pods nginx-1880671902-s3fdq

#### 删除deployment
*  kubectl delete deployment nginx

#### 更新Deployment
* kubectl set image deployment/nginx-deployment nginx=nginx:1.9.1

#### Deployment 扩容
* kubectl scale deployment nginx-deployment --replicas 10

#### 分析pod
*  kubectl describe pod nginx-1880671902-s3fdq

#### 看到应用日志
* kubectl logs nginx-dev-6c9744f59d-g92nl

#### Deployment方式 run 应用
* 创建 nginx-deployment.yaml
  ```
  apiVersion: apps/v1beta1
  kind: Deployment
  metadata:
    name: nginx-deployment
  spec:
    replicas: 3
    template:
      metadata:
        labels:
          app: nginx
      spec:
        containers:
        - name: nginx
          image: nginx:1.15
          ports:
          - containerPort: 80
  ```
* 运行: kubectl create -f nginx-deployment.yaml

* 删除: kubectl delete -f nginx-deployment.yaml


