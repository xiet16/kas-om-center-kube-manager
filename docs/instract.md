# k8s 组件架构

## 组件架构描述
1. master节点的基础组件包括 kube-ApiServer、kube-controller-manager、kube-scheduler
2. node节点组件包括 kubelet、kube-proxy、docker、pod。(https://blog.csdn.net/ambzheng/article/details/118877256)
3. pod 可以包含多个容器，容器间共享一个ip,栈等等
4. service 相当于网关层，用于同一服务在不同pod 间的负载均衡
5. kube-controller-manager 控制器管理包含 deployment(无状态控制器)、statefulSet(有状态控制器)、daemonSet(一次性控制器)、job(一次性任务)、cronJob(周期性任务)
6. 客户端使用kubectl 通过kube-ApiServer与k8s进行交互