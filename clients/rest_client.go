package clients

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	schema "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

/*
1、创建config，创建config的方法有两种：
(1) 使用clientcmd.BuildConfigFromFlags从配置文件中读取，读取的正是home目录下.kube/config中的内容。因此，需要将该文件拷贝到程序目录下。
(2) 使用rest.InClusterConfig()，该函数返回一个配置对象，它使用kubernetes提供给pods的服务账户。这个API是为运行在k8s的pod中的服务而设计的。如果你写的应用将来要跑在pod中，那么就可以使用这个方式。如果没有运行在kubernetes环境中的进程调用，它将会返回ErrNotInCluster。当我们的程序要在 pod中运行时，我们需要给其创建一个account，然后k8s会将需要的信息放入pod中，我们的程序就可以读取到了。
注意：使用restclient时还要自己设置config的GroupVersion和NegotiatedSerializer字段，如果不设置，会直接panic。同时也要设置APIPath，否则就查询不到 对应的资源。
2、使用config来创建restclient
3、使用restclient来对资源对象进行操作，比如查找指定命名空间下的pod等
*/

func Create() {
	//将.kube/config 下的配置文件拷贝到当前项目的config_files文件夹中,
	//如果两个字符串都为空字符串，那么将会尝试使用InClusterConfig来读取
	config, err := clientcmd.BuildConfigFromFlags("", "./config_files/config")
	if err != nil {
		panic(err)
	}
	config.GroupVersion = &coreV1.SchemeGroupVersion
	config.NegotiatedSerializer = schema.Codecs.WithoutConversion()
	config.APIPath = "/api"

	//构建rest client
	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	// 查找命名空间下的pod
	var podList coreV1.PodList
	if err = client.Get().Namespace("dev").Resource("pods").Do(context.Background()).Into(&podList); err != nil {
		log.Printf("get pod err:%v\n", err)
		return
	}
	for _, pod := range podList.Items {
		log.Printf("name: %s\n", pod.Name)
	}
}
