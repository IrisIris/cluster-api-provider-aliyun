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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ACKClusterSpec defines the desired state of ACKCluster
type ACKClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// cluster
	ClusterName        string `json:"cluster_name,omitempty"`
	ClusterType        string `json:"cluster_type,omitempty"`
	RegionId           string `json:"region_id,omitempty"`
	KubernetesVersion  string `json:"kubernetes_version"`
	CpuPolicy          string `json:"cpu_policy"`
	MasterInstanceType string `json:"master_instance_type"`
	NodesNum           int64  `json:"nodes_num"`
	WorkerInstanceType string `json:"master_instance_type"`

	// login
	LoginSpec LoginSpec `json:"login_spec"`
	// volume
	VolumeSpec VolumeSpec `json:"volume_spec"`
	// network
	NetworkSpec NetworkSpec `json:"networkSpec,omitempty"`
	Addons      Addons      `json:"addons"`

	Tags Tags `json:"tags"`
}
type LoginSpec struct {
	KeyPair       string `json:"key_pair"`
	LoginPassword string `json:"login_password"`
}
type NetworkSpec struct {
	// VPC ID，可空。如果不设置，系统会自动创建VPC，系统创建的VPC网段为192.168.0.0/16。
	// 说明 VpcId 和 vswitchid 只能同时为空或者同时都设置对应的值。
	VpcId            string   `json:"vpc_id""`
	MasterVswitchIds []string `json:"master_vswitch_ids"`
	WorkerVswitchds  []string `json:"worker_vswitchds"`

	SnatEntry *bool `json:"snat_entry"`
	// 容器网段，不能和VPC网段冲突。当选择系统自动创建VPC时，默认使用172.16.0.0/16网段。
	ContainerCidr string `json:"container_cidr"`
	// 服务网段，不能和VPC网段以及容器网段冲突。当选择系统自动创建VPC时，默认使用172.19.0.0/20网段。
	ServiceCidr string `json:"service_cidr"`

	EndpointPublicAccess *bool
}

type VolumeSpec struct {
	MasterSystemDisk ClusterSystemDisk `json:"master_system_disk"`
	WorkerSystemDisk ClusterSystemDisk `json:"master_system_disk"`
	DataDisk         []ClusterDataDisk `json:"data_disk"`
}

// system_disk
type ClusterSystemDisk struct {
	SystemDiskCategory string `json:"system_disk_category"`
	SystemDiskSize     string `json:"system_disk_size"`
}

// 挂载盘参数
type ClusterDataDisk struct {
	// category：数据盘类型。取值范围：
	//cloud：普通云盘
	//cloud_efficiency：高效云盘
	//cloud_ssd：SSD云盘
	Category string `json:"category"`

	Size int64

	Encrypted *bool
}

//
type Addons struct {
	Name    string
	Version string
	Config  string
}

// ACKClusterStatus defines the observed state of ACKCluster
type ACKClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//MasterIPs []string `json:"master_i_ps"`
	Ready             bool     `json:"ready"`
	MasterInstanceIDs []string `json:"master_instance_i_ds"`
	NodeInstanceIDs   []string `json:"node_instance_i_ds"`
	ScalingGroupID    string   `json:"scaling_group_id"`
	// 专有网络
	VpcId string `json:"vpc_id"`
	// 虚拟交换机
	VSwitchIds    string `json:"v_switch_ids"`
	IntranetSlbId string `json:"intranet_slb_id"`
	// ProxyMode:ipvs/iptables:"The mode we use in kube-proxy."
}

// +kubebuilder:object:root=true

// ACKCluster is the Schema for the ackclusters API
type ACKCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ACKClusterSpec   `json:"spec,omitempty"`
	Status ACKClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ACKClusterList contains a list of ACKCluster
type ACKClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ACKCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ACKCluster{}, &ACKClusterList{})
}
