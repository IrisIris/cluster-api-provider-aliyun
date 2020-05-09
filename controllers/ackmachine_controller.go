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
	"github.com/IrisIris/cluster-api-provider-aliyun/pkg/cloud/services"
	"github.com/IrisIris/cluster-api-provider-aliyun/pkg/cloud/services/ecs"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/record"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/IrisIris/cluster-api-provider-aliyun/api/v1alpha3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// ACKMachineReconciler reconciles a ACKMachine object
type ACKMachineReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	Recorder          record.EventRecorder
	ecsServiceFactory func() services.ECSMachineInterface
	//secretsManagerServiceFactory func(*scope.ClusterScope) services.SecretsManagerInterface
}

func (r *ACKMachineReconciler) getECSService() services.ECSMachineInterface {
	if r.ecsServiceFactory != nil {
		return r.ecsServiceFactory()
	}
	return ecs.NewService()
}

// +kubebuilder:rbac:groups=ack.cluster.k8s.io,resources=ackmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ack.cluster.k8s.io,resources=ackmachines/status,verbs=get;update;patch

func (r *ACKMachineReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("ackmachine", req.NamespacedName)

	// fetch the ACKMachine instance
	ackMachine := &infrav1.ACKMachine{}
	err := r.Get(ctx, req.NamespacedName, ackMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	machine, err := ecs.GetOwnerMachine(ackMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		logger.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	// fetch the cluster-api Cluster
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		logger.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}
	// whether ackMachine or cluster is marked as paused
	if isPaused(ackMachine) || isPaused(cluster) {
		logger.Info("ACKMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	// fetch aws cluster

	// Handle deleted machines

	// Handle not-deleted machines
	return r.reconcileNormals()

	return ctrl.Result{}, nil
}

func (r *ACKMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.ACKMachine{}).
		Complete(r)
}

func (r *ACKMachineReconciler) reconcileDeleted() {

}

func (r *ACKMachineReconciler) reconcileNormals(machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	machineScope.Info("Reconciling ACKMachine")

	// whether failed already
	if machineScope.HasFailed() {
		return ctrl.Result{}, nil
	}
	// add default Finalizer if not exits
	controllerutil.AddFinalizer(machineScope.ACKMachine, infrav1.MachineFinalizer)
	// todo {Register the finalizer immediately to avoid orphaning ACK resources on delete}
	if err := machineScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	// if cluster not ready
	if !machineScope.Cluster.Status.InfrastructureReady {
		machineScope.Info("Cluster is not ready yet")
		return ctrl.Result{}, nil
	}

	// todo {why not judge Bootstrap.Data}
	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Machine bootstrap data secret reference is not yet available")
		return ctrl.Result{}, nil
	}

	ecsSvc := r.getECSService()

	// get or create ecs instance
	instance, err := r.getOrCreate(machineScope, &ecsSvc)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Set an failure message if we couldn't find the instance.
	if instance == nil {
		machineScope.Info("EC2 instance cannot be found")
		machineScope.SetFailureReason(capierrors.UpdateMachineError)
		machineScope.SetFailureMessage(errors.New("EC2 instance cannot be found"))
		return ctrl.Result{}, nil
	}

	// Make sure Spec.ProviderID is always set.
	machineScope.SetProviderID(fmt.Sprintf("ack:////%s", instance.ID))

	existingInstanceState := machineScope.GetInstanceState()
	machineScope.SetInstanceState(instance.State)

	// Proceed to reconcile the AckMachine state.
	if existingInstanceState == nil || *existingInstanceState != instance.State {
		machineScope.Info("ECS instance state changed", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
	}

	// according to instance state to update ackMachine Status
	switch instance.State {
	case infrav1.InstanceStatePending, infrav1.InstanceStateStopping, infrav1.InstanceStateStopped:
		machineScope.SetNotReady()
	case infrav1.InstanceStateRunning:
		machineScope.SetReady()
	case infrav1.InstanceStateShuttingDown, infrav1.InstanceStateTerminated:
		machineScope.SetNotReady()
		machineScope.Info("Unexpected EC2 instance termination", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
		r.Recorder.Eventf(machineScope.ACKMachine, corev1.EventTypeWarning, "InstanceUnexpectedTermination", "Unexpected EC2 instance termination")
	default:
		machineScope.SetNotReady()
		machineScope.Info("EC2 instance state is undefined", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
		r.Recorder.Eventf(machineScope.ACKMachine, corev1.EventTypeWarning, "InstanceUnhandledState", "EC2 instance state is undefined")
		machineScope.SetFailureReason(capierrors.UpdateMachineError)
		machineScope.SetFailureMessage(errors.Errorf("EC2 instance state %q is undefined", instance.State))
	}

	// tasks that can take place during all known instance states, e.g. ensure tags

	// tasks that can only take place during operational instance states
	// e.g. slb, security groups
	return ctrl.Result{}, nil
}

func (r *ACKMachineReconciler) getOrCreate(scope *scope.MachineScope, ecsSvc *services.ECSMachineInterface) (*infrav1.Instance, error) {
	// first to get
	findOne, err := r.findInstance()
	if err != nil {

	}
	if findOne != nil {
		return findOne, nil
	}

	// Otherwise then create one
	// get userData
	userData, err := GetUserData()
	if err != nil {
		return nil, err
	}

	instance, err := (*ecsSvc).RunInstances()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create AWSMachine instance")
	}
	return instance, nil
}

func (r *ACKMachineReconciler) findInstance() (*infrav1.Instance, error) {
	// query ecs service to find instance
	return nil, nil
}
