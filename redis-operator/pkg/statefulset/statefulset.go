package statefulset

import (
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	datav1beta1 "redis-operator/api/v1beta1"
)

func NewStatefulSet(redis *datav1beta1.Redis) *v12.StatefulSet {
	labels := map[string]string{
		"app": redis.Name,
	}

	return &v12.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(redis, schema.GroupVersionKind{
					Group:   v12.SchemeGroupVersion.Group,
					Version: v12.SchemeGroupVersion.Version,
					Kind:    "AppService",
				}),
			},
		},
		Spec: v12.StatefulSetSpec{
			Replicas: redis.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: newContainers(redis),
				},

			},
			//ServiceName: "",
		},
	}
}

func newContainers(redis *datav1beta1.Redis) []v1.Container {
	return []v1.Container{
		{
			Name:  redis.Name,
			Image: redis.Spec.Image,
			Ports: []v1.ContainerPort{
				{
					ContainerPort: 6379,
				},
			},
			ImagePullPolicy: v1.PullIfNotPresent,
		},
	}
}
