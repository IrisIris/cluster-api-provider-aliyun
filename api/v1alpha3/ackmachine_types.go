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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/errors"
)

const MachineFinalizer = "ackmachine.infrastructure.cluster.x-k8s.io"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ACKMachineSpec defines the desired state of ACKMachine
// check for deatails: https://help.aliyun.com/document_detail/63440.html
type ACKMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ClusterId string `json:"cluster_id"`
	RegionId  string `json:"region_id,omitempty"`
	ZoneId    string `json:"zone_id"`
	// 实例的资源规格。如果您不指定LaunchTemplateId或LaunchTemplateName以确定启动模板，InstanceType为必选参数。
	InstanceType string `json:"instanceType"`
	InstanceName string `json:"instance_name"`
	Description  string `json:"description"`
	IoOptimized  string `json:"io_optimized"`

	ImageId  string   `json:"image_id"`
	Tags     Tags     `json:"tags"`
	UserData UserData `json:"user_data"`

	MachineNetworkSpec MachineNetworkSpec `json:"machine_network_spec"`
	MachineVolumeSpec  MachineVolumeSpec  `json:"machine_volume_spec"`

	// charge related
	// 是否要自动续费。当参数InstanceChargeType取值PrePaid时才生效。
	// 取值范围：true：自动续费。false（默认）：不自动续费。
	AutoRenew bool `json:"auto_renew"`
	// 每次自动续费的时长，当参数AutoRenew取值True时为必填。
	// PeriodUnit为Week时，AutoRenewPeriod取值{"1", "2", "3"}。
	//PeriodUnit为Month时，AutoRenewPeriod取值{"1", "2", "3", "6", "12"}。
	AutoRenewPeriod int64 `json:"auto_renew_period"`
}

type CpuOptions struct {
	// CPU核心数。该参数不支持自定义设置，只能采用默认值。
	Core           int64
	ThreadsPerCore int64
}

type MachineNetworkSpec struct {
	// SecurityGroupId决定了实例的网络类型，例如，如果指定安全组的网络类型为专有网络VPC，实例则为VPC类型，并同时需要指定参数VSwitchId。
	//如果您不指定LaunchTemplateId或LaunchTemplateName以确定启动模板，SecurityGroupId为必选参数。
	SecurityGroupId string `json:"security_group_id"`
	VSwitchId       string `json:"vswitch_id"`
	// 公网入带宽最大值，单位为Mbit/s。
	InternetMaxBandwidthIn  int64 `json:"internet_max_bandwidth_in"`
	InternetMaxBandwidthOut int64 `json:"internet_max_bandwidth_out"`
	// 网络计费类型
	// 取值范围：PayByBandwidth：按固定带宽计费。 PayByTraffic（默认）：按使用流量计费。
	InternetChargeType string `json:"internet_charge_type"`
	PrivateIpAddress   string `json:"private_ip_address"`
}

type MachineVolumeSpec struct {
	SystemDisk SystemDisk  `json:"system_disk"`
	DataDisks  []*DataDisk `json:"data_disks"`
	// 公网入带宽最大值，单位为Mbit/s。
	InternetMaxBandwidthIn  int64 `json:"internet_max_bandwidth_in"`
	InternetMaxBandwidthOut int64 `json:"internet_max_bandwidth_out"`
}

type UserData struct {
	Encryped *bool  `json:"encryped"`
	Datas    string `json:"datas""`
}

// ACKMachineStatus defines the observed state of ACKMachine
type ACKMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	MachineId      string                     `json:"machine_id"`
	InstanceId     string                     `json:"instance_id"`
	FailureReason  *errors.MachineStatusError `json:"failureReason,omitempty"`
	FailureMessage *string                    `json:"failureMessage,omitempty"`
}
type Tags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// todo
type NetworkInterface struct {
	// 辅助弹性网卡名称
	NetworkInterfaceName string
	Description          string
	//添加一张辅助弹性网卡并设置主IP地址。N的取值范围为1。
	//说明 创建ECS实例时，您最多能添加一张辅助网卡。实例创建成功后，您可以调用CreateNetworkInterface和AttachNetworkInterface添加更多的辅助网卡。
	//默认值：从网卡所属的交换机网段中随机选择一个IP地址。
	PrimaryIpAddress string
	//辅助弹性网卡所属的虚拟交换机ID。N的取值范围为1。
	//默认值：ECS实例所属的虚拟交换机。
	VSwitchId string
	// 互斥
	SecurityGroupId  string
	SecurityGroupIds []string
}
type SystemDisk struct {
	Size                 string `json:"size"`
	Category             string `json:"category"`
	DiskName             string `json:"disk_name"`
	Description          string
	PerformanceLevel     string
	AutoSnapshotPolicyId string
}

type DataDisk struct {
	Size                 string `json:"size"`
	SnapshotId           string `json:"snapshot_id"`
	Category             string `json:"category"`
	Encrypted            *bool  `json:"encrypted"`
	KMSKeyId             string `json:"kms_key_id"`
	DiskName             string `json:"disk_name"`
	Description          string `json:"description"`
	DeleteWithInstance   *bool  `json:"delete_with_instance"`
	PerformanceLevel     string `json:"performance_level"`
	AutoSnapshotPolicyId string `json:"auto_snapshot_policy_id"`
}

// +kubebuilder:object:root=true

// ACKMachine is the Schema for the ackmachines API
type ACKMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ACKMachineSpec   `json:"spec,omitempty"`
	Status ACKMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ACKMachineList contains a list of ACKMachine
type ACKMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ACKMachine `json:"items"`
}

func (am *ACKMachine) HasFailed() bool {
	return am.Status.FailureReason != nil || am.Status.FailureMessage != nil
}

func init() {
	SchemeBuilder.Register(&ACKMachine{}, &ACKMachineList{})
}
