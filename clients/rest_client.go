package clients

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	schema "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func Create() {
	//将.kube/config 下的配置文件拷贝到当前项目的config_files文件夹中
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
