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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	drzhanggv1beta1 "http-operator/api/v1beta1"
)

// FrigateReconciler reconciles a Frigate object
type FrigateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=drzhangg.my.domain,resources=frigates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=drzhangg.my.domain,resources=frigates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=drzhangg.my.domain,resources=frigates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Frigate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *FrigateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.WithValues("http ", req.NamespacedName)

	// TODO(user): your logic here
	var frigate drzhanggv1beta1.Frigate
	err := r.Get(ctx, req.NamespacedName, &frigate)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	gd := get()
	pd := post()

	fmt.Println("get:", string(gd))
	fmt.Println("post:", string(pd))

	logger.Info("fetch http operator objects", "httpservice", frigate)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FrigateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&drzhanggv1beta1.Frigate{}).
		Complete(r)
}

func get() []byte {
	resp, err := http.Get("http://127.0.0.1:32225/hello/tmd")
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}
	return body
}

func post() []byte {
	client := &http.Client{}

	url := "http://127.0.0.1:32225/user/drzhangg"
	m := map[string]interface{}{
		"age":     23,
		"gender":  1,
		"address": "上海沙田公寓24",
	}

	d, _ := json.Marshal(&m)
	req, err := http.NewRequest("POST", url, bytes.NewReader(d))
	if err != nil {
		// handle error
		fmt.Println("err:", err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("err:", err)
		return nil
	}

	return body
}
