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
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"redis-operator/pkg/statefulset"

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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=update;list;get;create;delete

func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.TODO()

	// v1:单机版
	logger := log.FromContext(ctx)
	logger.WithValues("redis", req.NamespacedName)

	// TODO(user): your logic here
	var redis datav1beta1.Redis

	// 先查询是否创建过 redis 实例
	if err := r.Client.Get(ctx, req.NamespacedName, &redis); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 有redis实例，对statefulSet进行创建或者更新
	var sts v1.StatefulSet
	if err := r.Client.Get(ctx, req.NamespacedName, &sts); err != nil && errors.IsNotFound(err) {
		nsts := statefulset.NewStatefulSet(&redis)
		if err := r.Client.Create(ctx, nsts); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datav1beta1.Redis{}).
		Complete(r)
}
