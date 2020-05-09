package scope

import (
	"context"
	providerv1 "github.com/IrisIris/cluster-api-provider-aliyun/api/v1alpha3"
	"github.com/IrisIris/cluster-api-provider-aliyun/pkg/cloud/services/ecs"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ClusterFinalizer allows ReconcileAWSCluster to clean up ACK resources associated with AWSCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "ackcluster.infrastructure.cluster.x-k8s.io"
)

type ClusterScopeParams struct {
	ACKClients
	Client     client.Client
	Logger     logr.Logger
	Cluster    *clusterv1.Cluster
	ACKCluster *providerv1.ACKCluster
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	ACKClients
	Cluster    *clusterv1.Cluster
	ACKCluster *providerv1.ACKCluster
}

func NewClusterScope(params *ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}

	if params.ACKCluster == nil {
		return nil, errors.New("failed to generate new scope from nil ACKCluster")
	}
	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	// TODO
	if params.ACKClients.ECS == nil {
		// new ecs client
	}

	helper, err := patch.NewHelper(params.ACKCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ClusterScope{
		Logger:      params.Logger,
		client:      params.Client,
		patchHelper: helper,
		ACKClients:  params.ACKClients,
		Cluster:     params.Cluster,
		ACKCluster:  params.ACKCluster,
	}, nil
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.PatchObject()
}

// PatchObject persists the cluster configuration and status.
func (s *ClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.ACKCluster)
}

func (s *ClusterScope) ReconcileDelete() (ctrl.Result, error) {
	s.Info("")
	ecsSvc := ecs.NewService()

	// todo delete network
	// todo delete load balancer

	// if cluster is deleted remove the finalizer
	controllerutil.RemoveFinalizer(s.ACKCluster, ClusterFinalizer)
	return ctrl.Result{}, nil
}

func (s *ClusterScope) ReconcileNormal() (reconcile.Result, error) {
	s.Info("Reconciling ACKCluster")
	ackCluster := s.ACKCluster

	// add finalizer if not exits
	controllerutil.AddFinalizer(ackCluster, ClusterFinalizer)
	if err := s.PatchObject(); err != nil {
		return reconcile.Result{}, nil
	}

	//ecsService := ecs.NewService()
	// todo:ReconcileNetwork
	//ecsService.ReconcileNetwork()
	// todo: ReconcileLoadbalancers

	ackCluster.Status.Ready = true
	return reconcile.Result{}, nil
}
