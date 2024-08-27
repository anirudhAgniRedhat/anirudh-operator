/*
Copyright 2024.

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

package controller

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	anirudhgroupv1 "github.com/anirudhAgniRedhat/anirudh-operator/api/v1"
	myv1 "github.com/anirudhAgniRedhat/anirudh-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MyCRDReconciler reconciles a MyCRD object
type MyCRDReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=anirudhgroup.k8s.io,resources=mycrds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=anirudhgroup.k8s.io,resources=mycrds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=anirudhgroup.k8s.io,resources=mycrds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyCRD object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MyCRDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mycrd", req.NamespacedName)

	// Fetch the MyCRD instance
	var myCRD myv1.MyCRD
	err := r.Get(ctx, req.NamespacedName, &myCRD)
	if err != nil {
		if errors.IsNotFound(err) {
			// CRD is deleted, recreate it
			log.Info("MyCRD resource not found. Recreating.")
			// Define a new MyCRD object and create it
			newMyCRD := &myv1.MyCRD{
				ObjectMeta: metav1.ObjectMeta{
					Name:      req.Name,
					Namespace: req.Namespace,
				},
				Spec: myv1.MyCRDSpec{
					// Your spec fields go here
				},
			}
			err = r.Create(ctx, newMyCRD)
			if err != nil {
				log.Error(err, "Failed to recreate MyCRD")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		// Error reading the object
		return ctrl.Result{}, err
	}

	// If the CRD exists, just log it
	log.Info("MyCRD resource found.")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyCRDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log = ctrl.Log.WithName("controllers").WithName("MyCRD")
	return ctrl.NewControllerManagedBy(mgr).
		For(&anirudhgroupv1.MyCRD{}).
		Complete(r)
}
