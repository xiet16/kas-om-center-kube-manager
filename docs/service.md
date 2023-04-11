# 说明
1. k8s 如果需要使用本地镜像，需要加载（minikube image load bc-order），且代码里设置ImagePullPolicy: corev1.PullNever
2. service 与pod 的关联通过 select labels, 查看关联关系用命令 kubectl get ep
3. 查看描述用命令 kubectl describe svc/pods svc-name/pod-name
4. 查看pod的标签用命令 kubectl get pods --show-labels
5. service 和pod 是一对多关系，并不会直接耦合
6. service 的三种方式及原理
7. minikube service service-name --url 可以开启minikube 的监听
8. ingress 待续


## 相关url
1. https://blog.csdn.net/LeoForBest/article/details/126879401
2. https://blog.csdn.net/huahua1999/article/details/124236875
