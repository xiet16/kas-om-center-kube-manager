package client_example

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	typeCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type ServiceManager struct {
	Client typeCoreV1.ServiceInterface
}

func NewServiceManager(ns string) *ServiceManager {
	if ns == "" {
		ns = corev1.NamespaceDefault
	}
	return &ServiceManager{
		Client: createClient(ns),
	}
}

func createClient(ns string) typeCoreV1.ServiceInterface {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatalln(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	deployClient := clientSet.CoreV1().Services(ns)
	return deployClient
}

func (sm *ServiceManager) CreateService(serverName string) {
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("service-%s", serverName),
			Labels: map[string]string{
				"name": serverName,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: "NodePort",
			Ports: []corev1.ServicePort{
				{
					Port:       serverPort,
					TargetPort: intstr.IntOrString{Type: 0, IntVal: serverPort},
					NodePort:   30085,
				},
			},
			Selector: map[string]string{
				"app": serverName,
			},
		},
	}
	if _, err := sm.Client.Create(context.Background(), &service, metav1.CreateOptions{}); err != nil {
		log.Fatalln(err)
	}
	log.Println("create service suc")
}
