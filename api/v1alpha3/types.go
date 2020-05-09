package v1alpha3

type Instance struct {
	Id           string `json:"id,omitempty"`
	InstanceName string `json:"instance_name"`

	State string `json:"state"`

	RegionId string `json:"region_id"`
	ZoneId   string `json:"zone_id"`

	//
	ResourceGroupId string `json:"resource_group_id"`
	// 实例规格 ecs.g5.large
	InstanceType string `json:"instance_type"`
	//实例规格族 ecs.g5
	InstanceTypeFamily string `json:"instance_type_family"`

	ImageId string `json:"image_id"`

	// The name of the SSH key pair.
	KeyPairName string `json:"key_pair_name"`

	SecurityGroupIDs []string `json:"security_group_i_ds,omitempty"`

	//
	CPU    int64 `json:"cpu"`
	Memory int64 `json:"memory"`
	// 实例的操作系统类型，分为Windows Server和Linux两种。可能值： windows linux
	OSType string `json:"os_type"`
	// 实例的操作系统名称 CentOS 7.4 64 位
	OSName string `json:"os_name"`
	// 实例操作系统的英文名称 CentOS 7.4 64 bit
	OSNameEn string `json:"os_name_en"`
	// UserData is the raw data script passed to the instance which is run upon bootstrap.
	// This field must not be base64 encoded and should only be used when running a new instance.
	UserData *string `json:"userData,omitempty"`

	// Network classic/vpc
	VlanId              string     `json:"vlan_id"`
	InstanceNetworkType string     `json:"instance_network_type"`
	InnerIpAddress      []string   `json:"inner_ip_address"`
	PublicIpAddress     []string   `json:"public_ip_address"`
	EipAddress          EipAddress `json:"eip_address"`

	// volume related
	DeviceAvailable *bool `json:"device_available"`

	// 实例的计费方式。可能值：PrePaid：包年包月。PostPaid：按量付费。
	InstanceChargeType string `json:"instance_charge_type"`

	//
	Tags []*Tag `json:"tags"`
}

type Tag struct {
	TagKey   string `json:"tag_key"`
	TagValue string `json:"tag_value"`
}
type VpcAttributes struct {
	// 云产品的IP，用于VPC云产品之间的网络互通。
	NatIpAddress string `json:"nat_ip_address"`
	// 私有IP地址
	PrivateIpAddress []string `json:"private_ip_address"`
	// 虚拟交换机ID
	VSwitchId string `json:"v_switch_id"`
	// 专有网络VPC ID
	VpcId string `json:"vpc_id"`
}

type InstanceState struct {
}

type EipAddress struct {
}
