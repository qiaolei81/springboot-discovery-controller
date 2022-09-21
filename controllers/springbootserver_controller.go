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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	springbootsitesv1 "microsoft.com/springboot-discovery-controller/api/v1"
)

// SpringBootServerReconciler reconciles a SpringBootServer object
type SpringBootServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=springbootsites.microsoft.com,resources=springbootservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=springbootsites.microsoft.com,resources=springbootservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=springbootsites.microsoft.com,resources=springbootservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpringBootServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *SpringBootServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	springBoot := &springbootsitesv1.SpringBootServer{}
	springBoot.Name = req.Name
	springBoot.Namespace = req.Namespace
	var err error
	var toDelete bool
	if err = r.Client.Get(ctx, req.NamespacedName, springBoot); err != nil {
		if errors.IsNotFound(err) {
			toDelete = true
		} else {
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	hooks := Hooks[springbootsitesv1.SpringBootServer]{
		Upsert: func(springBoot *springbootsitesv1.SpringBootServer) error {
			return r.Update(ctx, springBoot)
		},
		UpdateStatus: func(springBoot *springbootsitesv1.SpringBootServer) error {
			return r.Status().Update(ctx, springBoot)
		},
		Delete: func(discovered *springbootsitesv1.SpringBootServer) error {
			// Do nothing, we don't do delete SpringBootServer in controller
			return nil
		},
	}
	hooksForDiscovered := Hooks[springbootsitesv1.SpringBootDiscovered]{
		Upsert: func(discovered *springbootsitesv1.SpringBootDiscovered) error {
			existed := springbootsitesv1.SpringBootDiscovered{}
			if e := r.Get(ctx, req.NamespacedName, &existed); e != nil {
				if errors.IsNotFound(e) {
					return r.Create(ctx, discovered)
				} else {
					return err
				}
			}
			return r.Update(ctx, discovered)
		},
		UpdateStatus: func(discovered *springbootsitesv1.SpringBootDiscovered) error {
			return r.Status().Update(ctx, discovered)
		},
		Delete: func(discovered *springbootsitesv1.SpringBootDiscovered) error {
			return r.Delete(ctx, discovered)
		},
	}

	model := NewModel(springBoot, hooks, hooksForDiscovered, toDelete, log.Log)

	var e error
	err = model.Reconcile()
	if err != nil {
		e = model.markAsFailed(err)
	} else {
		e = model.markAsReconciled()
	}
	if e != nil {
		err = e
	}
	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpringBootServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&springbootsitesv1.SpringBootServer{}).
		Complete(r)
}
