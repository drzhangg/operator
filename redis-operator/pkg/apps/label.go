package apps

import "redis-operator/api/v1beta1"

func newStatefulSetLabel(cr *v1beta1.Redis) map[string]string {
	return map[string]string{
		"app": cr.Name,
	}
}
