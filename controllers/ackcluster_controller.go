/*
Copyright 2020 ALIYUN.

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
	"github.com/IrisIris/cluster-api-provider-aliyun/pkg/cloud/scope"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ackv1alpha3 "github.com/IrisIris/cluster-api-provider-aliyun/api/v1alpha3"
)

// ACKClusterReconciler reconciles a ACKCluster object
type ACKClusterReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=ack.cluster.k8s.io,resources=ackclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ack.cluster.k8s.io,resources=ackclusters/status,verbs=get;update;patch

func (r *ACKClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var reterr error
	ctx := context.Background()
	logger := r.Log.WithValues("ackcluster", req.NamespacedName, "ackCluster", req.Name)

	// fetch the ACKCluster instance
	ackCluster := &ackv1alpha3.ACKCluster{}
	err := r.Get(ctx, req.NamespacedName, ackCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// fetch the Cluster (cluster-api)
	cluster, err := util.GetOwnerCluster(ctx, r.Client, ackCluster.ObjectMeta)
	if err != nil {
		logger.Error(err, "failed to fetch the Cluster")
		return ctrl.Result{}, nil
	}
	if cluster == nil {
		logger.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if util.IsPaused(cluster, ackCluster) {
		logger.Info("ACKluster or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("cluster", cluster.Name)

	// create the scope
	clusterScope, err := scope.NewClusterScope(&scope.ClusterScopeParams{
		Client:     r.Client,
		Logger:     logger,
		Cluster:    cluster,
		ACKCluster: ackCluster,
	})
	if err != nil {
		return ctrl.Result{}, errors.Errorf("failed to create cluster scope: %+v", err)
	}

	defer func() {
		if err := clusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if r.IsDeletedACKCluster(ackCluster) {
		// todo: reconcile deleted
		_, err = clusterScope.ReconcileDelete()
		if err != nil {

		}
		//return reconcileDelete(clusterScope)
	}
	_, err = clusterScope.ReconcileNormal()
	if err != nil {

	}

	return ctrl.Result{}, nil
}

func (r *ACKClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ackv1alpha3.ACKCluster{}).
		Complete(r)
}

func (r *ACKClusterReconciler) IsDeletedACKCluster(ackCluster *ackv1alpha3.ACKCluster) bool {
	return !ackCluster.DeletionTimestamp.IsZero()
}
