package apps

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"redis-operator/api/v1beta1"
)

func newRedis(cr *v1beta1.Redis) *v1beta1.Redis {
	return &v1beta1.Redis{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    newStatefulSetLabel(cr),
		},
		Spec: v1beta1.RedisSpec{
			Image:    cr.Spec.Image,
			Replicas: cr.Spec.Replicas,
		},
	}
}
