package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"nginx-operator/api/v1beta1"
)

func NewDeploy(app *v1beta1.AppService) *appsv1.Deployment {

	labels := map[string]string{
		"app": app.Name,
	}
	selector := &metav1.LabelSelector{MatchLabels: labels}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "apps/v1",
			APIVersion: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   v1beta1.GroupVersion.Group,
					Version: v1beta1.GroupVersion.Version,
					Kind:    app.Kind,
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Size,
			Selector: selector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: newContainers(app),
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}
}

func newContainers(app *v1beta1.AppService) []v1.Container {
	containerPorts := []v1.ContainerPort{}

	for _, svcPorts := range app.Spec.Ports {
		cport := v1.ContainerPort{}
		cport.ContainerPort = svcPorts.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}

	return []v1.Container{
		{
			Name:            app.Name,
			Image:           app.Spec.Image,
			Resources:       app.Spec.Resource,
			Ports:           containerPorts,
			ImagePullPolicy: v1.PullIfNotPresent,
			Env:             app.Spec.Envs,
		},
	}
}
