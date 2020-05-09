package scope

import (
	"context"
	infrav1 "github.com/IrisIris/cluster-api-provider-aliyun/api/v1alpha3"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/klog/klogr"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineScopeParams struct {
	Client     client.Client
	Logger     logr.Logger
	Cluster    *clusterv1.Cluster
	Machine    *clusterv1.Machine
	ACKCluster *infrav1.ACKCluster
	ACKMachine *infrav1.ACKMachine
}
type MachineScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	Cluster    *clusterv1.Cluster
	Machine    *clusterv1.Machine
	ACKCkuster *infrav1.ACKCluster
	ACKMachine *infrav1.ACKMachine
}

func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	if params.Client == nil {
		return nil, errors.Errorf("failed to create machine scope due to empty client")
	}
	if params.Cluster == nil {
		return nil, errors.Errorf("failed to create machine scope due to empty cluster")
	}
	if params.Machine == nil {
		return nil, errors.Errorf("failed to create machine scope due to empty client")
	}
	if params.ACKCluster == nil {
		return nil, errors.Errorf("failed to create machine scope due to empty ACKCluster")
	}
	if params.ACKMachine == nil {
		return nil, errors.Errorf("failed to create machine scope due to empty ACKMachine")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	helper, err := patch.NewHelper(params.ACKMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &MachineScope{
		Logger:      params.Logger,
		client:      params.Client,
		patchHelper: helper,
		Cluster:     params.Cluster,
		Machine:     params.Machine,
		ACKCkuster:  params.ACKCluster,
		ACKMachine:  params.ACKMachine,
	}, nil
}

func (m *MachineScope) HasFailed() bool {
	return m.ACKMachine.Status.FailureReason != nil || m.ACKMachine.Status.FailureMessage != nil
}

// PatchObject persists the machine spec and status.
func (m *MachineScope) PatchObject() error {
	return m.patchHelper.Patch(context.TODO(), m.ACKMachine)
}

// SetFailureMessage sets the AWSMachine status failure message.
func (m *MachineScope) SetFailureMessage(v error) {
	m.ACKMachine.Status.FailureMessage = pointer.StringPtr(v.Error())
}

// SetFailureReason sets the AWSMachine status failure reason.
func (m *MachineScope) SetFailureReason(v capierrors.MachineStatusError) {
	m.ACKMachine.Status.FailureReason = &v
}
