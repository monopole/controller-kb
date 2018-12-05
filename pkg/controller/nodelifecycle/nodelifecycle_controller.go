/*

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

package nodelifecycle

import (
	"context"
	"log"

	bootyv1beta1 "github.com/monopole/controller-kb/pkg/apis/booty/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NodeLifeCycle Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNodeLifeCycle{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nodelifecycle-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to NodeLifeCycle
	err = c.Watch(&source.Kind{Type: &bootyv1beta1.NodeLifeCycle{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by NodeLifeCycle - change this for objects you create
	//err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &bootyv1beta1.NodeLifeCycle{},
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNodeLifeCycle{}

// ReconcileNodeLifeCycle reconciles a NodeLifeCycle object
type ReconcileNodeLifeCycle struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NodeLifeCycle object and makes changes based on the state read
// and what is in the NodeLifeCycle.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=booty.controller-kb.com,resources=nodelifecycles,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileNodeLifeCycle) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the NodeLifeCycle nodeLifeCycle
	log.Printf("Entering reconcile.")

	nodeLifeCycle := &bootyv1beta1.NodeLifeCycle{}
	err := r.Get(context.TODO(), request.NamespacedName, nodeLifeCycle)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if nodeLifeCycle.Spec.State != "reboot-desired" {
		log.Printf("No reboot desired!  Bailing out. %v", nodeLifeCycle.Spec.State)
		return reconcile.Result{}, nil
	}

	nodeLifeCycle.Spec.State = "reboot-NOW"
	log.Printf("Signalling reboot NOW!")

	err = r.Update(context.TODO(), nodeLifeCycle)
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
