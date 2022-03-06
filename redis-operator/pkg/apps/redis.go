package apps

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"redis-operator/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func newRedis(req ctrl.Request) *v1beta1.Redis {
	return &v1beta1.Redis{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ,
			Namespace: "",
			Labels:    nil,
		},
		Spec:   v1beta1.RedisSpec{},
		Status: v1beta1.RedisStatus{},
	}
}
