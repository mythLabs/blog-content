/*
Copyright 2025.

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
    "fmt"
    "time"

    "github.com/go-resty/resty/v2"
    "gopkg.in/yaml.v2"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    appsv1 "k8s.io/api/apps/v1"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
    ctrl "sigs.k8s.io/controller-runtime"
    
    deployv1alpha1 "github.com/mythLabs/blog-content/kubernetes-operator/deploysyncer-operator/api/v1alpha1"
)


// DeploySyncerReconciler reconciles a DeploySyncer object
type DeploySyncerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=deploy.example.com,resources=deploysyncers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=deploy.example.com,resources=deploysyncers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=deploy.example.com,resources=deploysyncers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeploySyncer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.4/pkg/reconcile
func (r *DeploySyncerReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the DeploySyncer CR
	deploySyncer := &v1.DeploySyncer{}
	err := r.Client.Get(ctx, req.NamespacedName, deploySyncer)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	// Fetch the deployments.yaml file from GitHub
	deploymentURL := fmt.Sprintf("%s/kubernetes-operator-app/deployments.yaml", deploySyncer.Spec.RepoURL, deploySyncer.Spec.Branch)
	client := resty.New()
	resp, err := client.R().Get(deploymentURL)
	if err != nil {
		log.Error(err, "Failed to fetch deployments.yaml")
		return reconcile.Result{}, err
	}

	if resp.StatusCode() != 200 {
		log.Error(err, "Failed to fetch deployments.yaml, status code", resp.Status())
		return reconcile.Result{}, fmt.Errorf("Failed to fetch deployments.yaml")
	}

	// Unmarshal the YAML into a Deployment object
	deployment := &appsv1.Deployment{}
	err = yaml.Unmarshal(resp.Body(), deployment)
	if err != nil {
		log.Error(err, "Failed to unmarshal deployments.yaml")
		return reconcile.Result{}, err
	}

	// Ensure the deployment is applied
	err = r.Client.Create(ctx, deployment)
	if err != nil && !controllerutil.IsAlreadyExists(err) {
		log.Error(err, "Failed to create deployment")
		return reconcile.Result{}, err
	}

	// Update DeploySyncer status
	deploySyncer.Status.LastAppliedTime = metav1.Now()
	deploySyncer.Status.LastStatus = "Deployment applied successfully"
	err = r.Client.Status().Update(ctx, deploySyncer)
	if err != nil {
		log.Error(err, "Failed to update DeploySyncer status")
		return reconcile.Result{}, err
	}

	log.Info("Deployment applied successfully")

	// Requeue after CheckInterval
	return reconcile.Result{
		RequeueAfter: time.Duration(deploySyncer.Spec.CheckInterval) * time.Second,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploySyncerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&deployv1alpha1.DeploySyncer{}).
		Named("deploysyncer").
		Complete(r)
}
