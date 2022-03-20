/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	datav1beta1 "redis-operator/api/v1beta1"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=data.my.domain,resources=redis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=data.my.domain,resources=redis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=data.my.domain,resources=redis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Redis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.TODO()

	// v1:单机版
	logger := log.FromContext(ctx)
	logger.WithValues("redis", req.NamespacedName)

	// TODO(user): your logic here
	var redis datav1beta1.Redis

	// 先查询是否创建过 redis 实例
	if err := r.Get(ctx, req.NamespacedName, &redis); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 有redis实例，对statefulset进行创建或者更新
	var sts v1.StatefulSet
	or, err := ctrl.CreateOrUpdate(ctx, r.Client, &sts, func() error {
		labels := map[string]string{
			"app": redis.Name,
		}
		sts.Spec = v1.StatefulSetSpec{
			Replicas: redis.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: v12.PodSpec{
					Containers: []v12.Container{
						{
							Name:       "redis",
							Image:      redis.Spec.Image,
							Command:    nil,
							Args:       nil,
							WorkingDir: "",
							Ports: []v12.ContainerPort{
								v12.ContainerPort{
									Name:          redis.Name,
									ContainerPort: 6379,
								},
							},
						},
					},
					RestartPolicy: v12.RestartPolicyAlways,
				},
			},
			ServiceName: "",
		}

		return ctrl.SetControllerReference(&redis, &sts, &runtime.Scheme{})
	})
	if err != nil {

		return ctrl.Result{}, err
	}
	logger.Info("CreateOrUpdate", "StatefulSet", or)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datav1beta1.Redis{}).
		Complete(r)
}
