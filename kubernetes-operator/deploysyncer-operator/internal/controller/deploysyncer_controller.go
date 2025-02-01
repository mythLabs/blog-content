package controller

import (
    "context"
    "fmt"
    "time"

    "github.com/go-resty/resty/v2"
    appsv1 "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"
    "sigs.k8s.io/yaml"

    deploysyncerv1alpha1 "github.com/mythLabs/blog-content/tree/main/kubernetes-operator/deploysyncer-operator/api/v1alpha1"
)

type DeploySyncerReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *DeploySyncerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)
    
    deploySyncer := &deploysyncerv1alpha1.DeploySyncer{}
    if err := r.Get(ctx, req.NamespacedName, deploySyncer); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Validate fields
    if deploySyncer.Spec.RepoURL == "" || deploySyncer.Spec.Branch == "" {
        deploySyncer.Status.LastStatus = "Missing required fields"
        r.Status().Update(ctx, deploySyncer)
        return ctrl.Result{}, fmt.Errorf("required fields missing")
    }

    // Setup GitHub client
    client := resty.New()
    deploymentURL := fmt.Sprintf("%s/%s/kubernetes-operator-app/deployments.yaml", 
        deploySyncer.Spec.RepoURL,
        deploySyncer.Spec.Branch)

    // Fetch deployment
    resp, err := client.R().Get(deploymentURL)
    if err != nil {
        deploySyncer.Status.LastStatus = fmt.Sprintf("Failed to fetch: %v", err)
        r.Status().Update(ctx, deploySyncer)
        return ctrl.Result{RequeueAfter: time.Minute}, err
    }

    // Parse deployment
    deployment := &appsv1.Deployment{}
    if err := yaml.Unmarshal(resp.Body(), deployment); err != nil {
        deploySyncer.Status.LastStatus = "Invalid YAML"
        r.Status().Update(ctx, deploySyncer)
        return ctrl.Result{}, err
    }

    // Apply deployment
    if err := r.Create(ctx, deployment); err != nil {
        if client.IgnoreAlreadyExists(err) != nil {
            if err := r.Update(ctx, deployment); err != nil {
                deploySyncer.Status.LastStatus = "Failed to update deployment"
                r.Status().Update(ctx, deploySyncer)
                return ctrl.Result{}, err
            }
        }
    }

    // Update status
    deploySyncer.Status.LastStatus = "Success"
    deploySyncer.Status.LastSyncTime = time.Now().Format(time.RFC3339)
    if err := r.Status().Update(ctx, deploySyncer); err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{RequeueAfter: time.Duration(deploySyncer.Spec.IntervalSeconds) * time.Second}, nil
}

func (r *DeploySyncerReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&deploysyncerv1alpha1.DeploySyncer{}).
        Owns(&appsv1.Deployment{}).
        Complete(r)
}