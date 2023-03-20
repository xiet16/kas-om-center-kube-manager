package client_example

import (
	"context"
	appsV1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsResV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"log"
)

type DeploymentManager struct {
	Client appsResV1.DeploymentInterface
}

func NewDeploymentManager() *DeploymentManager {
	return &DeploymentManager{}
}

func (dm *DeploymentManager) createClient() {
	config, err := clientcmd.BuildConfigFromFlags("", "./config_files")
	if err != nil {
		log.Fatalln(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	deployClient := clientSet.AppsV1().Deployments(corev1.NamespaceDefault)
	dm.Client = deployClient
}

func (dm *DeploymentManager) Create() {
	if dm.Client == nil {
		dm.createClient()
	}
	relicas := int32(2)
	deploy := appsV1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deploy-nginx-demo",
			Namespace: corev1.NamespaceDefault,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: &relicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nginx",
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []corev1.ContainerPort{
								{
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	dt, err := dm.Client.Create(context.Background(), &deploy, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("create deployment suc :", dt.Name)
}

func (dm *DeploymentManager) Update() {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deploy, err := dm.Client.Get(context.Background(), "deploy-nginx-demo", metav1.GetOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		replicas := int32(3)
		deploy.Spec.Replicas = &replicas
		deploy.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
		if _, err := dm.Client.Update(context.Background(), deploy, metav1.UpdateOptions{}); err != nil {
			log.Fatalln(err)
		}
		return err
	})
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func (dm *DeploymentManager) ListDeploy() {
	klog.Info("ListDeploy...........")
	deplist, err := dm.Client.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("list deployment error, err:%v", err)
		return
	}

	for _, dep := range deplist.Items {
		klog.Infof("deploy name:%s, replicas:%d, container image:%s", dep.Name, *dep.Spec.Replicas, dep.Spec.Template.Spec.Containers[0].Image)
	}
}

func (dm *DeploymentManager) DeleteDeploy() {
	klog.Info("DeleteDeploy...........")
	// 删除策略
	deletePolicy := metav1.DeletePropagationForeground
	err := dm.Client.Delete(context.Background(), "deploy-nginx-demo", metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		klog.Errorf("delete deployment error, err:%v", err)
	} else {
		klog.Info("delete deployment success")
	}
}
