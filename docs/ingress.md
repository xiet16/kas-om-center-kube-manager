# ingress
1. ingress 主要包含ingress-controller、nginx、ingress 
2. ingress 用于配置路由规则以及绑定service,可以理解为相当于nginx免去了手动reload
3. ingress-controller 通过api-server 监听资源的变化，从而更新ingress
4. ingress 使用例子可查看官网
5. 一个ingress 绑定多少service 合理
6. 如何查看ingress 的日志