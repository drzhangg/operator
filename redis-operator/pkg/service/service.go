package service

import (
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"redis-operator/api/v1beta1"
)

func NewService(redis *v1beta1.Redis) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(redis,schema.GroupVersionKind{
					Group:   v12.SchemeGroupVersion.Group,
					Version: v12.SchemeGroupVersion.Version,
					Kind:    "AppService",
				}),
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name: redis.Name,
					Port: 6379,
				},
			},
			Selector: map[string]string{
				"app": redis.Name,
			},
			Type: v1.ServiceTypeNodePort,
		},
	}
}
