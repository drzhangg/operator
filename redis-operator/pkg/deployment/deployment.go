package deployment

import (
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"redis-operator/api/v1beta1"
)

func NewDeployment(redis *v1beta1.Redis) *v1.Deployment {
	labels := map[string]string{
		"app": redis.Name,
	}

	return &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: redis.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					OwnerReferences: []metav1.OwnerReference{},
				},
				Spec: v12.PodSpec{
					Containers: newContainers(redis),
				},
			},
		},
	}
}

func newContainers(redis *v1beta1.Redis) []v12.Container {
	return []v12.Container{
		{
			Name:  redis.Name,
			Image: redis.Spec.Image,
			Ports: []v12.ContainerPort{
				{
					ContainerPort: 6379,
				},
			},
			ImagePullPolicy: v12.PullIfNotPresent,
		},
	}
}

