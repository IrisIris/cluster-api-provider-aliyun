package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type ECSMachineInterface interface {
	RunInstances(request *ecs.RunInstancesRequest) (response *ecs.RunInstancesResponse, err error)
}
