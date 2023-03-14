package clients

import (
	"bytes"
	"context"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"text/template"
)

func CreatePod() {
	//加载配置
	config, err := clientcmd.BuildConfigFromFlags("", "./conf/config")
	if err != nil {
		panic(err)
	}

	//创建clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//对资源对象进行操作
	podList, err := clientSet.CoreV1().Pods("dev").List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		log.Printf("list pods error:%v\n", err)
		return
	}
	for _, pod := range podList.Items {
		log.Printf("name: %s\n", pod.Name)
	}
}

func CreateNginxPod() {
	config, err := clientcmd.BuildConfigFromFlags("", "./config_files/config")
	if err != nil {
		panic(err)
	}

	//
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//初始化pod信息
	pod := coreV1.Pod{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
			Labels: map[string]string{
				"run": "nginx",
			},
		},
		Spec: coreV1.PodSpec{
			Containers: []coreV1.Container{
				{
					Image: "nginx",
					Name:  "nginx-container",
					Ports: []coreV1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	//创建pod
	_, err = clientSet.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	if err != nil {
		log.Printf("create pod error:%v\n", err)
		return
	}
	log.Printf("create pod success\n")
}

func CreatePodByTemplate() {
	config, err := clientcmd.BuildConfigFromFlags("", "./config_files/config")
	if err != nil {
		panic(err)
	}

	//创建clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//从模板文件构建pod
	spec := PodSpec{
		Name:          "nginx-pod-demo",
		Image:         "nginx",
		Namespace:     "default",
		ContainerName: "nginx",
	}

	var pod coreV1.Pod
	tmp1, err := ParseTemplate("./tmplate/pod.yaml", &spec)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(tmp1, &pod); err != nil {
		panic(err)
	}

	if _, err = clientSet.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{}); err != nil {
		log.Printf("create pod error:%v\n", err)
		return
	}
	log.Printf("create pod success\n")
}

type PodSpec struct {
	Name          string `json:"name"`
	Image         string `json:"image"`
	Namespace     string `json:"namespace"`
	ContainerName string `json:"container_name"`
}

func ParseTemplate(name string, item *PodSpec) ([]byte, error) {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, item)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
