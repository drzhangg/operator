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
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	tgv1beta1 "learn-operator/api/v1beta1"
)

// LearnReconciler reconciles a Learn object
type LearnReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tg.my.domain,resources=learns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tg.my.domain,resources=learns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tg.my.domain,resources=learns/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Learn object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *LearnReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	learn := tgv1beta1.Learn{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, &learn)
	if err != nil {
		return ctrl.Result{}, err
	}

	containerPorts := []v12.ContainerPort{}
	for _, v := range learn.Spec.Ports {
		cport := v12.ContainerPort{}
		cport.ContainerPort = v.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}

	deploy := v1.Deployment{}
	if err := r.Client.Get(context.TODO(), req.NamespacedName, &deploy); err != nil && errors.IsNotFound(err) {

		if err := r.Client.Create(context.TODO(), &v1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      learn.Name,
				Namespace: learn.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(&learn, schema.GroupVersionKind{
						Group:   tgv1beta1.GroupVersion.Group,
						Version: tgv1beta1.GroupVersion.Version,
						Kind:    "Learn",
					}),
				},
			},
			Spec: v1.DeploymentSpec{
				Replicas: learn.Spec.Replicas,
				Template: v12.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": learn.Name},
					},
					Spec: v12.PodSpec{
						Containers: []v12.Container{
							{
								Name:  learn.Name,
								Image: learn.Spec.Image,
								Ports: containerPorts,
							},
						},
					},
				},
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": learn.Name},
				},
			},
		}); err != nil {
			log.Log.Info("create deploy err:", err)
			return ctrl.Result{}, err
		}

		// create service

		if err := r.Client.Create(context.TODO(), &v12.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      learn.Name,
				Namespace: learn.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(&learn, schema.GroupVersionKind{
						Group:   tgv1beta1.GroupVersion.Group,
						Version: tgv1beta1.GroupVersion.Version,
						Kind:    "Learn",
					}),
				},
			},
			Spec: v12.ServiceSpec{
				Ports:    learn.Spec.Ports,
				Selector: map[string]string{"app": learn.Name},
				Type:     v12.ServiceTypeNodePort,
			},
		}); err != nil {
			log.Log.Info("create service err:", err)
			return ctrl.Result{}, err
		}
	}

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LearnReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tgv1beta1.Learn{}).
		Complete(r)
}
